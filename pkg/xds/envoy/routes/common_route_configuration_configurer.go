package routes

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	util_proto "github.com/kumahq/kuma/pkg/util/proto"
)

type CommonRouteConfigurationConfigurer struct {
	MaxDirectResponseBodySizeBytes uint32
}

func (c CommonRouteConfigurationConfigurer) Configure(routeConfiguration *envoy_config_route_v3.RouteConfiguration) error {
	if c.MaxDirectResponseBodySizeBytes == 0 {
		c.MaxDirectResponseBodySizeBytes = 4096
	}

	routeConfiguration.ValidateClusters = util_proto.Bool(false)
	routeConfiguration.MaxDirectResponseBodySizeBytes = util_proto.UInt32(c.MaxDirectResponseBodySizeBytes)
	return nil
}
