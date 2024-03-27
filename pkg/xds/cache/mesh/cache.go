package mesh

import (
	"context"
	"os"
	"time"

	"github.com/kumahq/kuma/pkg/core"
	"github.com/kumahq/kuma/pkg/metrics"
	"github.com/kumahq/kuma/pkg/xds/cache/once"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
)

// Cache is needed to share and cache Hashes among goroutines which
// reconcile Dataplane's state. Calculating hash is a heavy operation
// that requires fetching all the resources belonging to the Mesh.
type Cache struct {
	// cache is used for caching a context and ignoring mesh changes for up to a
	// short expiration time.
	cache *once.Cache
	// hashCache keeps a cached context, for a much longer time, that is only reused
	// when the mesh hasn't changed.
	hashCache xds_context.CustomCache

	meshContextBuilder xds_context.MeshContextBuilder

	usingRedisCache bool
}

// cleanupTime is the time after which the mesh context is removed from
// the longer TTL cache.
// It exists to ensure contexts of deleted Meshes are eventually cleaned up.
const cleanupTime = 15 * time.Minute

var logger = core.Log.WithName("xds").WithName("cache")

func NewCache(
	expirationTime time.Duration,
	meshContextBuilder xds_context.MeshContextBuilder,
	metrics metrics.Metrics,
	redisUrl string,
) (*Cache, error) {
	c, err := once.New(expirationTime, "mesh_cache", metrics)
	if err != nil {
		return nil, err
	}

	usingRedisCache := false
	var meshCtxCache xds_context.CustomCache
	meshCtxCache = xds_context.NewInMemoryCache(cleanupTime)
	if os.Getenv("USE_REDIS_CACHE") != "" {
		logger.Info("Attempting to use redis cache for mesh contexts")

		meshCtxCache, err = xds_context.NewRedisCache(redisUrl, cleanupTime)
		if err != nil {
			logger.Error(err, "no redis cache setup for base mesh contexts")
			meshCtxCache = xds_context.NewInMemoryCache(cleanupTime)
		} else {
			usingRedisCache = true
		}
	}

	// cache.New(cleanupTime, time.Duration(int64(float64(cleanupTime)*0.9)))
	return &Cache{
		cache:              c,
		meshContextBuilder: meshContextBuilder,
		hashCache:          meshCtxCache,
		usingRedisCache:    usingRedisCache,
	}, nil
}

func (c *Cache) GetMeshContext(ctx context.Context, mesh string) (xds_context.MeshContext, error) {
	if c.usingRedisCache {
		return c.getMeshContextImpl(ctx, mesh)
	}

	// Check our short TTL cache for a context, ignoring whether there have been
	// changes since it was generated.
	elt, err := c.cache.GetOrRetrieve(ctx, mesh, once.RetrieverFunc(func(ctx context.Context, key string) (interface{}, error) {
		return c.getMeshContextImpl(ctx, mesh)
	}))
	if err != nil {
		return xds_context.MeshContext{}, err
	}
	return elt.(xds_context.MeshContext), nil
}

func (c *Cache) getMeshContextImpl(ctx context.Context, mesh string) (xds_context.MeshContext, error) {
	// Check hashCache first for an existing mesh latestContext
	latestContext := &xds_context.MeshContext{}
	_, err := c.hashCache.Get(ctx, mesh, latestContext)
	if err != nil {
		logger.Error(err, "could not fetch latest mesh context from redis cache", "mesh", mesh)
	}

	// Rebuild the context only if the hash has changed
	latestContext, err = c.meshContextBuilder.BuildIfChanged(ctx, mesh, latestContext)
	if err != nil {
		return xds_context.MeshContext{}, err
	}

	// By always setting the mesh context, we refresh the TTL
	// with the effect that often used contexts remain in the cache while no
	// longer used contexts are evicted.
	err = c.hashCache.Set(ctx, mesh, latestContext)
	if err != nil {
		logger.Error(err, "could not set latest mesh context in redis cache", "mesh", mesh)
	}

	return *latestContext, nil
}
