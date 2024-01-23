package routes

import (
	kumaroutes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
)

func CommonRouteConfiguration(maxDirectResponseBodySizeBytes uint32) kumaroutes.RouteConfigurationBuilderOpt {
	return kumaroutes.AddRouteConfigurationConfigurer(
		&kumaroutes.CommonRouteConfigurationConfigurer{
			MaxDirectResponseBodySizeBytes: maxDirectResponseBodySizeBytes,
		},
	)
}
