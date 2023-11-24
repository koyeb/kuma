package generator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	model "github.com/kumahq/kuma/pkg/core/xds"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_common "github.com/kumahq/kuma/pkg/xds/envoy"
	envoy_clusters "github.com/kumahq/kuma/pkg/xds/envoy/clusters"
	envoy_listeners "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
	"github.com/kumahq/kuma/pkg/xds/envoy/names"
	envoy_routes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
	envoy_virtual_hosts "github.com/kumahq/kuma/pkg/xds/envoy/virtualhosts"
)

const (
	// OriginProbes is a marker to indicate by which ProxyGenerator resources were generated.
	OriginProbe                      = "probe"
	listenerName                     = "probe:listener"
	routeConfigurationName           = "probe:route_configuration"
	TcpOriginProbe                   = "tcp_probe"
	tcpListenerNamePrefix            = "tcp_probe:listener:"
	tcpHealthCheckOriginProbe        = "tcp_health_check_probe"
	tcpHealthCheckServerProbeCluster = "tcp_health_check_server_probe"
	tcpHealthCheckNamePrefix         = tcpHealthCheckOriginProbe + ":listener:"
)

type ProbeProxyGenerator struct{}

func (g ProbeProxyGenerator) generateTcpProbes(ctx context.Context, xdsCtx xds_context.Context, proxy *model.Proxy) (*model.ResourceSet, error) {
	tcpProbes := proxy.Dataplane.Spec.TcpProbes
	tcpHealthCheckServerProbe := proxy.Dataplane.Spec.TcpHealthCheckServerProbe
	resources := model.NewResourceSet()
	localAddr := "127.0.0.1"

	if tcpHealthCheckServerProbe != nil {
		cluster, err := envoy_clusters.NewClusterBuilder(proxy.APIVersion, tcpHealthCheckServerProbeCluster).
			Configure(envoy_clusters.ProvidedEndpointCluster(
				false,
				core_xds.Endpoint{
					Target: localAddr,
					Port:   tcpHealthCheckServerProbe.DestinationPort,
				}),
			).
			Build()
		if err != nil {
			return nil, err
		}

		resources.Add(&core_xds.Resource{
			Name:     cluster.GetName(),
			Origin:   tcpHealthCheckServerProbeCluster,
			Resource: cluster,
		})

		tcpHealthCheckListenerName := fmt.Sprintf("%s%d", tcpHealthCheckNamePrefix, tcpHealthCheckServerProbe.Port)
		tcpHealthCheckListener, err := envoy_listeners.NewListenerBuilder(proxy.APIVersion, tcpHealthCheckListenerName).
			Configure(
				envoy_listeners.InboundListener(
					proxy.Dataplane.Spec.GetNetworking().GetAddress(),
					tcpHealthCheckServerProbe.Port,
					model.SocketAddressProtocolTCP,
				),
			).
			Configure(
				envoy_listeners.FilterChain(
					envoy_listeners.NewFilterChainBuilder(proxy.APIVersion, tcpHealthCheckListenerName).
						Configure(
							envoy_listeners.TcpProxyDeprecated(
								tcpHealthCheckListenerName,
								envoy_common.NewCluster(
									envoy_common.WithService(tcpHealthCheckServerProbeCluster),
								),
							),
						),
				),
			).
			Build()
		if err != nil {
			return nil, errors.Wrapf(err, "could not generate TCP health check listener %s", tcpHealthCheckListenerName)
		}
		resources.Add(&model.Resource{
			Name:     tcpHealthCheckListenerName,
			Resource: tcpHealthCheckListener,
			Origin:   tcpHealthCheckOriginProbe,
		})
	}

	for _, tcpProbe := range tcpProbes {
		tcpListenerName := fmt.Sprintf("%s%d", tcpListenerNamePrefix, tcpProbe.Port)
		tcpProbeListener, err := envoy_listeners.NewListenerBuilder(proxy.APIVersion, tcpListenerName).
			Configure(envoy_listeners.InboundListener(proxy.Dataplane.Spec.GetNetworking().GetAddress(), tcpProbe.Port, model.SocketAddressProtocolTCP)).
			Configure(envoy_listeners.FilterChain(
				envoy_listeners.NewFilterChainBuilder(proxy.APIVersion, tcpListenerName).
					Configure(
						envoy_listeners.TcpProxyDeprecated(
							tcpListenerName,
							envoy_common.NewCluster(
								envoy_common.WithService(names.GetLocalClusterName(tcpProbe.DestinationPort)),
							),
						),
					),
			)).
			Configure(envoy_listeners.OriginalDstForwarder()).
			Configure(envoy_listeners.TransparentProxying(proxy.Dataplane.Spec.Networking.GetTransparentProxying())).
			Build()
		if err != nil {
			return nil, errors.Wrapf(err, "could not generate listener %s", listenerName)
		}
		resources.Add(&model.Resource{
			Name:     tcpListenerName,
			Resource: tcpProbeListener,
			Origin:   TcpOriginProbe,
		})
	}

	return resources, nil

	return nil, nil
}

func (g ProbeProxyGenerator) Generate(ctx context.Context, _ *model.ResourceSet, xdsCtx xds_context.Context, proxy *model.Proxy) (*model.ResourceSet, error) {
	resources, err := g.generateTcpProbes(ctx, xdsCtx, proxy)
	if err != nil {
		return nil, err
	}

	probes := proxy.Dataplane.Spec.Probes
	if probes == nil {
		return resources, nil
	}

	virtualHostBuilder := envoy_virtual_hosts.NewVirtualHostBuilder(proxy.APIVersion, "probe")

	portSet := map[uint32]bool{}
	for _, inbound := range proxy.Dataplane.Spec.Networking.Inbound {
		portSet[proxy.Dataplane.Spec.Networking.ToInboundInterface(inbound).WorkloadPort] = true
	}
	for _, endpoint := range probes.Endpoints {
		matchURL, err := url.Parse(endpoint.Path)
		if err != nil {
			return nil, err
		}
		newURL, err := url.Parse(endpoint.InboundPath)
		if err != nil {
			return nil, err
		}
		if portSet[endpoint.InboundPort] {
			virtualHostBuilder.Configure(
				envoy_virtual_hosts.Route(matchURL.Path, newURL.Path, names.GetLocalClusterName(endpoint.InboundPort), true))
		} else {
			// On Kubernetes we are overriding probes for every container, but there is no guarantee that given
			// probe will have an equivalent in inbound interface (ex. sidecar that is not selected by any service).
			// In this situation there is no local cluster therefore we are sending redirect to a real destination.
			// System responsible for using virtual probes needs to support redirect (kubelet on K8S supports it).
			virtualHostBuilder.Configure(
				envoy_virtual_hosts.Redirect(matchURL.Path, newURL.Path, true, endpoint.InboundPort))
		}
	}

	probeListener, err := envoy_listeners.NewInboundListenerBuilder(proxy.APIVersion, proxy.Dataplane.Spec.GetNetworking().GetAddress(), probes.Port, model.SocketAddressProtocolTCP).
		WithOverwriteName(listenerName).
		Configure(envoy_listeners.FilterChain(envoy_listeners.NewFilterChainBuilder(proxy.APIVersion, envoy_common.AnonymousResource).
			Configure(envoy_listeners.HttpConnectionManager(listenerName, false)).
			Configure(envoy_listeners.HttpStaticRoute(envoy_routes.NewRouteConfigurationBuilder(proxy.APIVersion, routeConfigurationName).
				Configure(envoy_routes.VirtualHost(virtualHostBuilder)))))).
		Configure(envoy_listeners.TransparentProxying(proxy.Dataplane.Spec.Networking.GetTransparentProxying())).
		Build()
	if err != nil {
		return nil, errors.Wrapf(err, "could not generate listener %s", listenerName)
	}

	resources.Add(&model.Resource{
		Name:     listenerName,
		Resource: probeListener,
		Origin:   OriginProbe,
	})

	return resources, nil
}
