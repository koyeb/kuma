package sync

import (
	"context"
	"net"
	"sort"

	"github.com/koyeb/koyeb-api-client-go-internal/api/v1/koyeb"
	"github.com/kumahq/kuma/pkg/coord"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

type GlobalLoadBalancerProxyBuilder struct {
	*DataplaneProxyBuilder
	CatalogDatacenters  koyeb.CatalogDatacentersApi
	InternalDeployments koyeb.InternalDeploymentsApi
}

func (p *GlobalLoadBalancerProxyBuilder) Build(ctx context.Context, key core_model.ResourceKey, meshContext xds_context.MeshContext) (*core_xds.Proxy, error) {
	proxy, err := p.DataplaneProxyBuilder.Build(ctx, key, meshContext)
	if err != nil {
		return nil, err
	}

	datacenters, err := p.fetchDatacenters(ctx)
	if err != nil {
		return nil, err
	}

	koyebApps, err := p.fetchKoyebApps(ctx)
	if err != nil {
		return nil, err
	}

	endpointMap, err := p.buildEndpointMap(datacenters)
	if err != nil {
		return nil, err
	}

	proxy.GlobalLoadBalancerProxy = &core_xds.GlobalLoadBalancerProxy{
		Datacenters: datacenters,
		EndpointMap: endpointMap,
		KoyebApps:   koyebApps,
	}

	return proxy, nil
}

func (p *GlobalLoadBalancerProxyBuilder) fetchDatacenters(ctx context.Context) ([]*core_xds.KoyebDatacenter, error) {
	resp, _, err := p.CatalogDatacenters.ListDatacentersExecute(p.CatalogDatacenters.ListDatacenters(ctx))
	if err != nil {
		return nil, err
	}

	datacenters := []*core_xds.KoyebDatacenter{}
	for _, koyebDC := range *resp.Datacenters {
		coordinates, err := coord.NewCoord(*koyebDC.Coordinates)
		if err != nil {
			return nil, err
		}

		datacenters = append(datacenters, &core_xds.KoyebDatacenter{
			ID:       *koyebDC.Id,
			RegionID: *koyebDC.RegionId,
			Domain:   *koyebDC.Domain,
			Coord:    coordinates,
		})

	}

	return datacenters, nil
}

func (p *GlobalLoadBalancerProxyBuilder) fetchKoyebApps(ctx context.Context) ([]*core_xds.KoyebApp, error) {

	resp, _, err := p.InternalDeployments.ListAllRoutesExecute(p.InternalDeployments.ListAllRoutes(ctx).UseKumaV2(true))
	if err != nil {
		return nil, err
	}

	routesByAppId := map[string][]koyeb.ListAllRoutesReplyRoute{}

	// First, grroup routes by AppId
	for _, route := range resp.GetRoutes() {
		routesByAppId[route.GetAppId()] = append(routesByAppId[route.GetAppId()], route)
	}

	// Build output list
	koyebApps := []*core_xds.KoyebApp{}
	for _, routes := range routesByAppId {
		koyebApp := &core_xds.KoyebApp{}

		domains := map[string]struct{}{}
		for _, route := range routes {
			dcs := map[string]struct{}{}

			for _, dc := range route.GetDatacenters() {
				dcs[dc] = struct{}{}
			}

			domains[route.GetDomain()] = struct{}{}
			koyebApp.Services = append(koyebApp.Services, &core_xds.KoyebService{
				ID:              route.GetServiceId(),
				DatacenterIDs:   dcs,
				Port:            uint32(route.GetPort()),
				DeploymentGroup: route.GetDeploymentGroup(),
				Path:            route.GetPath(),
			})
		}

		koyebApp.Domains = maps.Keys(domains)
		koyebApps = append(koyebApps, koyebApp)
	}

	return koyebApps, nil
}

func (p *GlobalLoadBalancerProxyBuilder) buildEndpointMap(datacenters []*core_xds.KoyebDatacenter) (core_xds.EndpointMap, error) {
	endpointMap := core_xds.EndpointMap{}

	for _, dc := range datacenters {
		// TODO(nicoche): move that in a goroutine
		ips, err := net.LookupIP(dc.Domain)
		if err != nil {
			return nil, errors.Wrapf(err, "could not resolve %s", dc.Domain)
		}

		// Not sure that we need it, but always generate the slice of endpoints in
		// the same order.
		sort.Slice(ips, func(i, j int) bool {
			return ips[i].String() < ips[j].String()
		})

		for _, ip := range ips {
			endpointMap[dc.ID] = append(endpointMap[dc.ID], core_xds.Endpoint{
				Target: ip.String(),
				// TODO(nicoche) should we harmonize that? urgh it's ugly
				Port:   5602,
				Weight: 1,
			})
		}
	}

	return endpointMap, nil
}
