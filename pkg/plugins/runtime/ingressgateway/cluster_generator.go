package ingressgateway

import (
	"context"
	"fmt"
	"sort"

	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/metadata"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	"github.com/kumahq/kuma/pkg/xds/envoy/clusters"
	envoy_endpoints "github.com/kumahq/kuma/pkg/xds/envoy/endpoints"
	envoy_tags "github.com/kumahq/kuma/pkg/xds/envoy/tags"
	"github.com/kumahq/kuma/pkg/xds/generator/zoneproxy"
	"golang.org/x/exp/maps"
)

type ClusterGenerator struct {
}

func (c *ClusterGenerator) GenerateClusters(ctx context.Context, xdsCtx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, error) {
	resources := core_xds.NewResourceSet()

	availableSvcsByMesh := map[string][]*mesh_proto.ZoneIngress_AvailableService{}
	for _, service := range proxy.ZoneIngressProxy.ZoneIngressResource.Spec.AvailableServices {
		availableSvcsByMesh[service.Mesh] = append(availableSvcsByMesh[service.Mesh], service)
	}

	for _, mr := range proxy.ZoneIngressProxy.MeshResourceList {
		meshName := mr.Mesh.GetMeta().GetName()
		services := maps.Keys(mr.EndpointMap)
		sort.Strings(services)
		dest := zoneproxy.BuildMeshDestinations(
			availableSvcsByMesh[meshName],
			xds_context.Resources{MeshLocalResources: mr.Resources},
		)

		for _, service := range services {
			deploymentGroup := "prod"
			//deploymentGroup := dest[service]//["koyeb.com/deploymentgroup"]

			// CDS
			clusterName := fmt.Sprintf("%s_%s", service, deploymentGroup)

			r, err := generateEdsCluster(proxy, clusterName, service, dest)
			if err != nil {
				return nil, err
			}
			resources.Add(r)

			// EDS
			r, err = generateLoadAssignment(proxy, clusterName, mr.EndpointMap[service])
			if err != nil {
				return nil, err
			}
			resources.Add(r)
		}
	}

	return resources, nil
}

func generateEdsCluster(proxy *core_xds.Proxy, clusterName string, service string, dest map[string][]envoy_tags.Tags) (*core_xds.Resource, error) {
	tagSlice := envoy_tags.TagsSlice(append(dest[service], dest[mesh_proto.MatchAllTag]...))
	tagKeySlice := tagSlice.ToTagKeysSlice().Transform(
		envoy_tags.Without(mesh_proto.ServiceTag),
	)

	clusterBuilder := clusters.NewClusterBuilder(proxy.APIVersion, clusterName).Configure(
		clusters.EdsCluster(),
		clusters.LbSubset(tagKeySlice),
		//TODO(nicoche): MTLS
		clusters.ConnectionBufferLimit(DefaultConnectionBuffer),
	)

	r, err := buildClusterResource(clusterBuilder)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func buildClusterResource(c *clusters.ClusterBuilder) (*core_xds.Resource, error) {
	msg, err := c.Build()
	if err != nil {
		return nil, err
	}

	cluster := msg.(*envoy_cluster_v3.Cluster)

	return &core_xds.Resource{
		Name:     cluster.GetName(),
		Origin:   metadata.OriginIngressGateway,
		Resource: cluster,
	}, nil
}

func generateLoadAssignment(proxy *core_xds.Proxy, clusterName string, endpoints []core_xds.Endpoint) (*core_xds.Resource, error) {
	cla, err := envoy_endpoints.CreateClusterLoadAssignment(clusterName, endpoints, proxy.APIVersion)
	if err != nil {
		return nil, err
	}

	return &core_xds.Resource{
		Name:     clusterName,
		Origin:   metadata.OriginIngressGateway,
		Resource: cla,
	}, nil
}
