package context

import (
	"fmt"
	"context"
	"time"
	"golang.org/x/exp/maps"
	"encoding/base64"
	"slices"

	"github.com/kumahq/kuma/pkg/core"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/xds/cache/sha256"
)

var logger = core.Log.WithName("xds-server")

type meshContextFetcher = func(ctx context.Context, meshName string) (MeshContext, error)

func getReport(meshPerTiming map[time.Duration][]string) string {
	durations := maps.Keys(meshPerTiming)
	slices.Sort(durations)
	// ORder from slowest to fastest
	slices.Reverse(durations)

	distribution := []string{}
	for _, duration := range durations {
		distribution = append(distribution, fmt.Sprintf("%d (%s)", duration, meshPerTiming[duration]))
	}

	return fmt.Sprintf("Median duration: %d. Mean duration: %d. Distribution: %s.", distribution)
}

func AggregateMeshContexts(
	ctx context.Context,
	resManager manager.ReadOnlyResourceManager,
	fetcher meshContextFetcher,
) (AggregatedMeshContexts, error) {
	l := logger.WithName("aggregate-mesh-contexts")
	l.Info("AggregateMeshContexts(): listing meshes")

	var meshList core_mesh.MeshResourceList
	if err := resManager.List(ctx, &meshList, core_store.ListOrdered()); err != nil {
		return AggregatedMeshContexts{}, err
	}

	var meshContexts []MeshContext
	meshContextsByName := map[string]MeshContext{}
	l.Info("AggregateMeshContexts(): fetching mesh context for each mesh", "count_meshes", len(meshList.Items))
	meshPerTiming := map[time.Duration][]string{}
	for _, mesh := range meshList.Items {
		startAt := time.Now()
		meshCtx, err := fetcher(ctx, mesh.GetMeta().GetName())
		if err != nil {
			if core_store.IsResourceNotFound(err) {
				// When the mesh no longer exists it's likely because it was removed since, let's just skip it.
				continue
			}
			return AggregatedMeshContexts{}, err
		}
		meshContexts = append(meshContexts, meshCtx)
		meshContextsByName[mesh.Meta.GetName()] = meshCtx

		duration := time.Since(startAt)
		meshPerTiming[duration] = append(meshPerTiming[duration], mesh.Meta.GetName())
	}

	report := getReport(meshPerTiming)
	l.Info("AggregateMeshContexts(): mesh fetching report", "report", report)
	l.Info("AggregateMeshContexts(): hash mesh contexts")
	hash := aggregatedHash(meshContexts)

	l.Info("AggregateMeshContexts(): listing egress")
	egressByName := map[string]*core_mesh.ZoneEgressResource{}
	if len(meshContexts) > 0 {
		for _, egress := range meshContexts[0].Resources.ZoneEgresses().Items {
			egressByName[egress.Meta.GetName()] = egress
		}
	} else {
		var egressList core_mesh.ZoneEgressResourceList
		if err := resManager.List(ctx, &egressList, core_store.ListOrdered()); err != nil {
			return AggregatedMeshContexts{}, err
		}

		for _, egress := range egressList.GetItems() {
			egressByName[egress.GetMeta().GetName()] = egress.(*core_mesh.ZoneEgressResource)
		}

		hash = base64.StdEncoding.EncodeToString(core_model.ResourceListHash(&egressList))
	}

	result := AggregatedMeshContexts{
		Hash:               hash,
		Meshes:             meshList.Items,
		MeshContextsByName: meshContextsByName,
		ZoneEgressByName:   egressByName,
	}
	return result, nil
}

func aggregatedHash(meshContexts []MeshContext) string {
	var hash string
	for _, meshCtx := range meshContexts {
		hash += meshCtx.Hash
	}
	return sha256.Hash(hash)
}
