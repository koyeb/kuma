package context

import (
	"context"
	"time"

	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/log"
)

type AggregateMeshContextsComponent struct {
	resManager  manager.ReadOnlyResourceManager
	fetcher     meshContextFetcher
	parallelism uint32
}

func NewAggregateMeshContextsComponent(
	resManager manager.ReadOnlyResourceManager,
	fetcher meshContextFetcher,
	parallelism uint32,
) (*AggregateMeshContextsComponent, error) {

	return &AggregateMeshContextsComponent{
		resManager:  resManager,
		fetcher:     fetcher,
		parallelism: parallelism,
	}, nil
}

// We're running this routine every hour to warmup the cache for AggregateMeshContexts
func (a *AggregateMeshContextsComponent) Start(stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	l := log.AddFieldsFromCtx(logger, ctx, context.Background())

	_, err := a.Build(ctx)
	if err != nil {
		l.Error(err, "could not build aggregate mesh contexts")
	}
	l.Info("successfuly built aggregate mesh contexts")

	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-stop:
			return nil

		case <-ticker.C:
			// We don't record the result (yet?). We're mostly doing this to warm the
			// cache at startup.
			_, err := a.Build(ctx)
			if err != nil {
				l.Error(err, "could not build aggregate mesh contexts")
			}
			l.Info("successfuly built aggregate mesh contexts")
		}
	}

	return nil
}

func (a *AggregateMeshContextsComponent) NeedLeaderElection() bool {
	return false
}

func (a *AggregateMeshContextsComponent) Build(ctx context.Context) (AggregatedMeshContexts, error) {
	l := log.AddFieldsFromCtx(logger, ctx, context.Background())

	var meshList core_mesh.MeshResourceList
	if err := a.resManager.List(ctx, &meshList, core_store.ListOrdered()); err != nil {
		return AggregatedMeshContexts{}, err
	}

	meshes := make(chan string, len(meshList.Items))
	errors := make(chan error, len(meshList.Items))

	for i := 0; i < int(a.parallelism); i++ {
		go a.worker(ctx, meshes, errors)
	}

	for _, mesh := range meshList.Items {
		meshes <- mesh.GetMeta().GetName()
	}
	close(meshes)

	for i := 0; i < len(meshList.Items); i++ {
		err := <-errors
		if err != nil {
			l.Error(err, "could not build mesh context")
		}
	}

	return AggregatedMeshContexts{}, nil
}

func (a *AggregateMeshContextsComponent) worker(ctx context.Context, meshes <-chan string, errors chan<- error) {
	for meshName := range meshes {
		_, err := a.fetcher(ctx, meshName)
		errors <- err
	}
}
