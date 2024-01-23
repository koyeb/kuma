package ingressgateway

import (
	"fmt"
	"net/http"
	"sort"

	envoy_config_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/http_pages"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/route"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_tags "github.com/kumahq/kuma/pkg/xds/envoy/tags"
	"github.com/kumahq/kuma/pkg/xds/envoy/tls"
	"github.com/kumahq/kuma/pkg/xds/generator/zoneproxy"
	"golang.org/x/exp/maps"
)

func GenerateRouteBuilders(proxy *core_xds.Proxy) ([]*route.RouteBuilder, error) {
	routeBuilders := []*route.RouteBuilder{}

	// NOTE(nicoche)
	// This has been taken from pkg//xds/generator/zoneproxy/generator.go
	//
	// That's a lot of for loops. A few notes:
	//   * Some of the computation done here is also done in ClusterGenerator. We
	//     could factorize that part to reduce the number of iterations.
	//   * I'm not 100% sure about what happens from `for _, destination := ...`
	//   * The ZoneIngress proxy bases its routing on SNI. I guess that it allows
	//     it not to generate a config as big as our. We'd need to investigate
	//     that
	availableSvcsByMesh := map[string][]*mesh_proto.ZoneIngress_AvailableService{}
	for _, service := range proxy.ZoneIngressProxy.ZoneIngressResource.Spec.AvailableServices {
		availableSvcsByMesh[service.Mesh] = append(availableSvcsByMesh[service.Mesh], service)
	}

	for _, mr := range proxy.ZoneIngressProxy.MeshResourceList {
		meshName := mr.Mesh.GetMeta().GetName()
		services := maps.Keys(mr.EndpointMap)
		sort.Strings(services)

		sniUsed := map[string]struct{}{}
		availableServices := availableSvcsByMesh[meshName]
		destinationsPerService := zoneproxy.BuildMeshDestinations(
			availableSvcsByMesh[meshName],
			xds_context.Resources{MeshLocalResources: mr.Resources},
		)
		endpointMap := mr.EndpointMap

		for _, service := range availableServices {
			serviceName := service.Tags[mesh_proto.ServiceTag]
			deploymentGroup := service.Tags[mesh_proto.KoyebDeploymentGroup]

			destinations := destinationsPerService[serviceName]
			destinations = append(destinations, destinationsPerService[mesh_proto.MatchAllTag]...)

			// NOTE(nicoche): see if we should grab this dynamically
			serviceEndpoints := endpointMap[serviceName]

			relevantTagsUsed := map[string]struct{}{}

			for _, destination := range destinations {
				sni := tls.SNIFromTags(destination.
					WithTags(mesh_proto.ServiceTag, serviceName).
					WithTags("mesh", service.Mesh),
				)
				if _, ok := sniUsed[sni]; ok {
					continue
				}
				sniUsed[sni] = struct{}{}

				// relevantTags is a set of tags for which it actually makes sense to do LB split on.
				// If the endpoint list is the same with or without the tag, we should just not do the split.
				// However, we should preserve full SNI, because the client expects Zone Proxy to support it.
				// This solves the problem that Envoy deduplicate endpoints of the same address and different metadata.
				// example 1:
				// Ingress1 (10.0.0.1) supports service:a,version:1 and service:a,version:2
				// Ingress2 (10.0.0.2) supports service:a,version:1 and service:a,version:2
				// If we want to split by version, we don't need to do LB subset on version.
				//
				// example 2:
				// Ingress1 (10.0.0.1) supports service:a,version:1
				// Ingress2 (10.0.0.2) supports service:a,version:2
				// If we want to split by version, we need LB subset.
				relevantTags := envoy_tags.Tags{}
				for key, value := range destination {
					matchedTargets := map[string]struct{}{}
					allTargets := map[string]struct{}{}
					for _, endpoint := range serviceEndpoints {
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

				if _, ok := relevantTagsUsed[relevantTags.String()]; ok {
					continue
				}
				relevantTagsUsed[relevantTags.String()] = struct{}{}

				routeBuilder := &route.RouteBuilder{}
				routeBuilder.Configure(route.RouteMatchPrefixPath("/"))
				routeBuilder.Configure(route.RouteMatchHeaderExactMatch("X-KOYEB-ROUTE", fmt.Sprintf("%s_%s", serviceName, deploymentGroup)))
				routeBuilder.Configure(route.RouteActionClusterHeader("X-KOYEB-ROUTE", relevantTags, MaxStreamDuration))

				routeBuilders = append(routeBuilders, routeBuilder)
			}
		}
	}

	routeBuilders = append(routeBuilders, generateNoDestinationRouteBuilder())

	return routeBuilders, nil
}

// In case we receive a request with X-Koyeb-Route set, it means that the GLB
// let it go through. It means that the target service is starting, healthy
// or unhealthy.
//   * If the target service is healthy, a cluster will be defined and hence the
//     request will be caught by the block before this one
//   * If the target service is starting or unhealthy, the request will not be
//     caught until now. The following route builder will catch it and return a
//     custom 503 page.
//     We return a 503 because our CDN remaps 502s to custom pages that we
//     upload.
//
// Note that this route might catch some other requests which we do not want. For
// example, if Kuma considers the destination dataplanes as unhealthy (e.g.
// instance died unexpectedly), this will be matched. This is a known limitation
func generateNoDestinationRouteBuilder() *route.RouteBuilder {
	routeBuilder := &route.RouteBuilder{}

	routeBuilder.Configure(route.RouteMatchPrefixPath("/"))
	routeBuilder.Configure(route.RouteMatchPresentHeader("X-KOYEB-ROUTE", true))
	routeBuilder.Configure(route.RouteActionDirectResponse(http.StatusServiceUnavailable, http_pages.RouteNotPropagated))
	routeBuilder.Configure(route.RouteAddResponseHeader(&envoy_config_core.HeaderValueOption{Header: &envoy_config_core.HeaderValue{
		Key:   "Content-Type",
		Value: "text/html; charset=UTF-8",
	}}))

	return routeBuilder
}
