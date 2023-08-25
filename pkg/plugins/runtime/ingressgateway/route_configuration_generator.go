package ingressgateway

import (
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	envoy_names "github.com/kumahq/kuma/pkg/xds/envoy/names"
	envoy_routes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
)

func GenerateRouteConfig(protocol mesh_proto.MeshGateway_Listener_Protocol, proxy *core_xds.Proxy) *envoy_routes.RouteConfigurationBuilder {
	switch protocol {
	case mesh_proto.MeshGateway_Listener_HTTPS,
		mesh_proto.MeshGateway_Listener_HTTP:
	default:
		return nil
	}

	resourceName := envoy_names.GetGatewayListenerName("whatever", mesh_proto.MeshGateway_Listener_HTTP.String(), uint32(5601))
	return envoy_routes.NewRouteConfigurationBuilder(proxy.APIVersion, resourceName).
		Configure(
			envoy_routes.CommonRouteConfiguration(),
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
