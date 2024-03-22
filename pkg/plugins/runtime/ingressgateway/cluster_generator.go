package ingressgateway

import (
	"context"
	"fmt"
	"sort"
	"time"

	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/durationpb"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/metadata"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	"github.com/kumahq/kuma/pkg/xds/envoy/clusters"
	envoy_endpoints "github.com/kumahq/kuma/pkg/xds/envoy/endpoints"
	envoy_tags "github.com/kumahq/kuma/pkg/xds/envoy/tags"
	"github.com/kumahq/kuma/pkg/xds/generator/zoneproxy"
)

var ClusterCircuitBreakerSettings = core_mesh.CircuitBreakerResource{
	Spec: &mesh_proto.CircuitBreaker{
		Conf: &mesh_proto.CircuitBreaker_Conf{
			Thresholds: &mesh_proto.CircuitBreaker_Conf_Thresholds{
				MaxConnections:     util_proto.UInt32(1024),
				MaxPendingRequests: util_proto.UInt32(1024),
				MaxRequests:        util_proto.UInt32(1024),
				MaxRetries:         util_proto.UInt32(1024),
			},
			SplitExternalAndLocalErrors: true,
			BaseEjectionTime:            durationpb.New(120 * time.Second),
			MaxEjectionPercent:          &wrappers.UInt32Value{Value: 90},
			Detectors: &mesh_proto.CircuitBreaker_Conf_Detectors{
				GatewayErrors: &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{Consecutive: &wrappers.UInt32Value{Value: 3}},
				// make sure that it's smaller than the value of `MaxRetries` in retry policy
				LocalErrors: &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{Consecutive: &wrappers.UInt32Value{Value: 2}},
				Failure: &mesh_proto.CircuitBreaker_Conf_Detectors_Failure{
					RequestVolume: &wrappers.UInt32Value{Value: 5},
					MinimumHosts:  &wrappers.UInt32Value{Value: 1},
					Threshold:     &wrappers.UInt32Value{Value: 85},
				},
			},
		},
	},
}

type ClusterGenerator struct{}

func (c *ClusterGenerator) GenerateClusters(ctx context.Context, xdsCtx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, error) {
	resources := core_xds.NewResourceSet()

	availableSvcsByMesh := map[string][]*mesh_proto.ZoneIngress_AvailableService{}
	for _, service := range proxy.ZoneIngressProxy.ZoneIngressResource.Spec.AvailableServices {
		availableSvcsByMesh[service.Mesh] = append(availableSvcsByMesh[service.Mesh], service)
	}

	for _, mr := range proxy.ZoneIngressProxy.MeshResourceList {
		targetMesh := mr.Mesh
		targetMeshName := targetMesh.GetMeta().GetName()
		services := maps.Keys(mr.EndpointMap)
		sort.Strings(services)

		dest := zoneproxy.BuildMeshDestinations(
			availableSvcsByMesh[targetMeshName],
			xds_context.Resources{MeshLocalResources: mr.Resources},
		)

		for _, service := range services {
			// NOTE(nicoche): see if we should grab this dynamically
			clusterName := fmt.Sprintf("%s_%s", service, "prod")

			// CDS
			r, err := generateEdsCluster(proxy, clusterName, service, dest, xdsCtx.Mesh.Resource, targetMesh, mr.EndpointMap[service])
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

func getLbSplitTags(destinationTags []envoy_tags.Tags, endpoints []core_xds.Endpoint) envoy_tags.TagKeysSlice {
	out := envoy_tags.TagKeysSlice{}
	for _, destination := range destinationTags {
		relevantTags := envoy_tags.Tags{}
		for key, value := range destination {
			matchedTargets := map[string]struct{}{}
			allTargets := map[string]struct{}{}
			for _, endpoint := range endpoints {
				address := endpoint.Address()
				if endpoint.Tags[key] == value || value == mesh_proto.MatchAllTag {
					matchedTargets[address] = struct{}{}
				}
				allTargets[address] = struct{}{}
			}
			if len(matchedTargets) < len(allTargets) {
				relevantTags[key] = value
			}
		}

		if len(relevantTags) > 0 {
			out = append(out, relevantTags.Keys())
		}
	}

	return out
}

func generateEdsCluster(
	proxy *core_xds.Proxy,
	clusterName string,
	service string,
	dest map[string][]envoy_tags.Tags,
	sourceMesh *core_mesh.MeshResource,
	targetMesh *core_mesh.MeshResource,
	endpoints []core_xds.Endpoint,
) (*core_xds.Resource, error) {
	tagSlice := envoy_tags.TagsSlice(append(dest[service], dest[mesh_proto.MatchAllTag]...))
	lbSplitTagsKeys := getLbSplitTags(dest[service], endpoints)

	clusterBuilder := clusters.NewClusterBuilder(proxy.APIVersion, clusterName).Configure(
		clusters.EdsCluster(),
		clusters.LbSubset(lbSplitTagsKeys),
		clusters.CrossMeshClientSideMTLS(proxy.SecretsTracker, sourceMesh, targetMesh, service, true, tagSlice),
		clusters.ConnectionBufferLimit(DefaultConnectionBuffer),
		clusters.HttpDownstreamProtocolOptions(),
		clusters.CircuitBreaker(&ClusterCircuitBreakerSettings),
		clusters.OutlierDetection(&ClusterCircuitBreakerSettings),
		clusters.PanicMode(0),
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
