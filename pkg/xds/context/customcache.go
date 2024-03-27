package context

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CustomCache interface {
	Set(ctx context.Context, meshName string, item cachable) error
	Get(ctx context.Context, meshName string, receiver cachable) (bool, error)
}

func NewRedisCache(address string, defaultTTL time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := client.Get(ctx, "whatever").Result()
	if err != redis.Nil {
		return nil, errors.Wrapf(err, "could not contact redis at address %s", address)
	}

	return &RedisCache{
		client:     client,
		defaultTTL: defaultTTL,
	}, nil
}

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func (gc *GlobalContext) CacheKeyName(_ string) string {
	// shared between all meshes
	return "global-context"
}

func (gc *BaseMeshContext) CacheKeyName(meshName string) string {
	return fmt.Sprintf("base-context-mesh-%s", meshName)
}

type cachable interface {
	CacheKeyName(meshName string) string
}

func (rc *RedisCache) Set(ctx context.Context, meshName string, item cachable) error {
	if rc == nil {
		return nil
	}

	key := item.CacheKeyName(meshName)

	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(item)
	if err != nil {
		return err
	}

	err = rc.client.Set(ctx, key, b.Bytes(), rc.defaultTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) Get(ctx context.Context, meshName string, receiver cachable) (bool, error) {
	if rc == nil {
		return false, nil
	}

	key := receiver.CacheKeyName(meshName)

	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err = d.Decode(receiver)
	if err != nil {
		return false, err
	}

	return true, nil
}

type InMemoryCache struct {
	client *gocache.Cache
}

func NewInMemoryCache(defaultTTL time.Duration) *InMemoryCache {
	return &InMemoryCache{
		client: gocache.New(defaultTTL, time.Duration(int64(float64(defaultTTL)*0.9))),
	}
}

func (mc *InMemoryCache) Set(_ context.Context, meshName string, item cachable) error {
	if mc == nil || mc.client == nil {
		return nil
	}

	key := item.CacheKeyName(meshName)
	mc.client.SetDefault(key, item)

	return nil
}

func (mc *InMemoryCache) Get(_ context.Context, meshName string, receiver cachable) (bool, error) {
	if mc == nil || mc.client == nil {
		return false, nil
	}

	key := receiver.CacheKeyName(meshName)
	cached, ok := mc.client.Get(key)
	if !ok {
		return false, nil
	}

	switch v := receiver.(type) {
	case *GlobalContext:
		globalContext, ok := cached.(*GlobalContext)
		if !ok {
			return false, errors.New("cached value is not a *GlobalContext")
		}

		*v = *globalContext
		return true, nil

	case *BaseMeshContext:
		baseMeshContext, ok := cached.(*BaseMeshContext)
		if !ok {
			return false, errors.New("cached value is not a *BaseMeshContext")
		}

		*v = *baseMeshContext
		return true, nil
	}

	return false, errors.New("unhandled type, only handling *GlobalContext and *BaseMeshContext")
}
