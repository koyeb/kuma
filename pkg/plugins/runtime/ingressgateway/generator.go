package ingressgateway

import (
	"context"
	"fmt"

	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/pkg/errors"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/gateway/metadata"
	envoy_routes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_listeners "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
)

const (
	IngressGatewayRoutesName = "ingress-gateway-routes"
)

// FilterChainGenerator is responsible for handling the filter chain for
// a specific protocol.
// A FilterChainGenerator can be host-specific or shared amongst hosts.
type FilterChainGenerator interface {
	Generate(xdsCtx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, []*envoy_listeners.FilterChainBuilder, error)
}

// Generator generates xDS resources for an entire Ingress Gateway.
type Generator struct {
	Zone                     string
	HTTPFilterChainGenerator FilterChainGenerator
}

type Route struct {
	Mesh            string
	Service         string
	DeploymentGroup string
}

func (r Route) key() string {
	return r.DeploymentGroup
}

func (g Generator) Generate(
	ctx context.Context,
	xdsCtx xds_context.Context,
	proxy *core_xds.Proxy,
) (*core_xds.ResourceSet, error) {
	resources := core_xds.NewResourceSet()

	var limits []RuntimeResoureLimitListener

	// NOTE(nicoche) We're supposed to iterate on listeners here but for now we
	// only keep one listener
	ldsResources, limit, err := g.generateLDS(ctx, xdsCtx, proxy)
	if err != nil {
		return nil, err
	}
	resources.AddSet(ldsResources)

	if limit != nil {
		limits = append(limits, *limit)
	}

	rdsResources, err := g.generateRDS(mesh_proto.MeshGateway_Listener_HTTP, xdsCtx, proxy)
	if err != nil {
		return nil, err
	}
	resources.AddSet(rdsResources)

	// TODO(nicoche) end iteration supposed to be here

	resources.Add(g.generateRTDS(limits))

	return resources, nil
}

func (g Generator) generateRTDS(limits []RuntimeResoureLimitListener) *core_xds.Resource {
	layer := map[string]interface{}{}
	for _, limit := range limits {
		layer[fmt.Sprintf("envoy.resource_limits.listener.%s.connection_limit", limit.Name)] = limit.ConnectionLimit
	}

	res := &core_xds.Resource{
		Name:   "ingressgateway.listeners",
		Origin: metadata.OriginGateway,
		Resource: &envoy_service_runtime_v3.Runtime{
			Name:  "ingressgateway.listeners",
			Layer: util_proto.MustStruct(layer),
		},
	}

	return res
}

func (g Generator) generateLDS(ctx context.Context, xdsCtx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, *RuntimeResoureLimitListener, error) {
	resources := core_xds.NewResourceSet()

	listenerBuilder, limit := GenerateListener(proxy)

	//hostname := "*"

	res, filterChainBuilders, err := g.HTTPFilterChainGenerator.Generate(xdsCtx, proxy)
	if err != nil {
		return nil, limit, err
	}
	resources.AddSet(res)

	for _, filterChainBuilder := range filterChainBuilders {
		listenerBuilder.Configure(envoy_listeners.FilterChain(filterChainBuilder))
	}

	res, err = BuildResourceSet(listenerBuilder)
	if err != nil {
		return nil, limit, errors.Wrapf(err, "failed to build listener resource")
	}
	resources.AddSet(res)

	return resources, limit, nil
}

func (g Generator) generateRDS(protocol mesh_proto.MeshGateway_Listener_Protocol, xdsCtx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, error) {
	switch protocol {
	case mesh_proto.MeshGateway_Listener_HTTPS,
		mesh_proto.MeshGateway_Listener_HTTP:
	default:
		return nil, nil
	}

	resources := core_xds.NewResourceSet()
	routeConfig := GenerateRouteConfig(mesh_proto.MeshGateway_Listener_HTTP, proxy)

	vh, err := GenerateVirtualHost(xdsCtx, proxy, nil)
	if err != nil {
		return nil, err
	}

	routeConfig.Configure(envoy_routes.VirtualHost(vh))

	res, err := BuildResourceSet(routeConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build route configuration resource")
	}
	resources.AddSet(res)

	return resources, nil
}

