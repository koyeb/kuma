package ingressgateway

import (
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/route"
)

func GenerateRouteBuilders(proxy *core_xds.Proxy) ([]*route.RouteBuilder, error) {
	routeBuilders := []*route.RouteBuilder{}

	routeBuilder := &route.RouteBuilder{}
	routeBuilder.Configure(route.RouteMatchPrefixPath("/"))
	routeBuilder.Configure(route.RouteMatchPresentHeader("X-KOYEB-ROUTE", true))
	routeBuilder.Configure(route.RouteActionClusterHeader("X-KOYEB-ROUTE"))

	routeBuilders = append(routeBuilders, routeBuilder)

	return routeBuilders, nil
}
