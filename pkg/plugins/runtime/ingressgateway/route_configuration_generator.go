package ingressgateway

import (
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/routes"
	envoy_routes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
)

func GenerateRouteConfig(info GatewayListenerInfo) *envoy_routes.RouteConfigurationBuilder {
	switch info.Listener.Protocol {
	case mesh_proto.MeshGateway_Listener_HTTPS,
		mesh_proto.MeshGateway_Listener_HTTP:
	default:
		return nil
	}

	return envoy_routes.NewRouteConfigurationBuilder(info.Proxy.APIVersion, info.Listener.ResourceName).
		Configure(
			//NOTE(nicoche): we have to expand this. The default is 4096 bytes. However, we embed
			// a full HTML page that we send over as a direct response in some cases.
			routes.CommonRouteConfiguration(uint32(32768)),
			envoy_routes.IgnorePortInHostMatching(),
			// TODO(jpeach) propagate merged listener tags.
			// Ideally we would propagate the tags header
			// to mesh services but not to external services,
			// but in the route configuration, we don't know
			// yet where the request will route to.
			// envoy_routes.TagsHeader(...),
			envoy_routes.ResetTagsHeader(),
		)

	// TODO(jpeach) apply additional route configuration configuration.
}
