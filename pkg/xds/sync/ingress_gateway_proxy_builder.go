package sync

import (
	"context"

	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	"github.com/kumahq/kuma/pkg/xds/envoy"
)

const (
	//TODO(nicoche) use something else
	zoneIngressName string = "zoneingress-par1"
)

type IngressGatewayProxyBuilder struct {
	*IngressProxyBuilder
}

func (p *IngressGatewayProxyBuilder) Build(
	ctx context.Context,
	key core_model.ResourceKey,
	aggregatedMeshCtxs xds_context.AggregatedMeshContexts,
) (*core_xds.Proxy, error) {
	meshContext := aggregatedMeshCtxs.MustGetMeshContext(key.Mesh)

	dp, found := meshContext.DataplanesByName[key.Name]
	if !found {
		return nil, core_store.ErrorResourceNotFound(core_mesh.DataplaneType, key.Name, key.Mesh)
	}

	// NOTE(nicoche)
	// We want to fetch a lot of context related to available endpoints
	// in the current zone, for all meshes.
	// However, Dataplanes are scoped to a single mesh. We hack around
	// by constructing a ZoneIngressProxy, which turns out to contain
	// pretty much all we need. This allows us to grab everything we
	// need without having to write a whole builder from scratch.

	// Here, reference an existing ZoneIngress to take it from DB and
	// update its AvailableServices field that is full of the info we
	// need to build the ZoneIngressProxy
	zoneIngressKey := core_model.ResourceKey{
		Name: zoneIngressName,
		Mesh: "", // A ZoneIngress is not part of any mesh
	}
	zoneIngress, err := p.getZoneIngress(ctx, zoneIngressKey, aggregatedMeshCtxs)
	if err != nil {
		return nil, err
	}

	meshName := meshContext.Resource.GetMeta().GetName()

	allMeshNames := []string{meshName}
	for _, mesh := range meshContext.Resources.OtherMeshes().Items {
		allMeshNames = append(allMeshNames, mesh.GetMeta().GetName())
	}

	secretsTracker := envoy.NewSecretsTracker(meshName, allMeshNames)

	proxy := &core_xds.Proxy{
		Id:               core_xds.FromResourceKey(key),
		APIVersion:       p.apiVersion,
		Dataplane:        dp,
		ZoneIngressProxy: p.buildZoneIngressProxy(zoneIngress, aggregatedMeshCtxs),
		SecretsTracker:   secretsTracker,
		Metadata:         &core_xds.DataplaneMetadata{},
	}
	return proxy, nil
}
