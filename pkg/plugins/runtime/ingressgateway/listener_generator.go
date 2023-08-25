package ingressgateway

import (
	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	envoy_listeners "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
	envoy_names "github.com/kumahq/kuma/pkg/xds/envoy/names"
)

// TODO(jpeach) It's a lot to ask operators to tune these defaults,
// and we probably would never do that. However, it would be convenient
// to be able to update them for performance testing and benchmarking,
// so at some point we should consider making these settings available,
// perhaps on the Gateway or on the Dataplane.

// Buffer defaults.
const DefaultConnectionBuffer = 32 * 1024

type RuntimeResoureLimitListener struct {
	Name            string
	ConnectionLimit uint32
}

// TODO(nicoche) take GatewayListenerInfo as input, maybe
func GenerateListener(proxy *core_xds.Proxy) (*envoy_listeners.ListenerBuilder, *RuntimeResoureLimitListener) {
	// TODO(nicoche) change, it comes from meshGateway

	var port uint32 = 5601

	protocol := mesh_proto.MeshGateway_Listener_HTTP
	address := proxy.Dataplane.Spec.GetNetworking().Address

	log.V(1).Info("generating listener",
		"address", address,
		"port", port,
		"protocol", protocol,
	)

	// TODO(nicoche): whatever -> from meshGateway
	name := envoy_names.GetGatewayListenerName("whatever", protocol.String(), port)
	// NOTE(nicoche_: Maybe use MeshGateway resources here instead of no limit
	var limits *RuntimeResoureLimitListener

	return envoy_listeners.NewInboundListenerBuilder(
		proxy.APIVersion,
		address,
		port,
		core_xds.SocketAddressProtocolTCP,
	).
		WithOverwriteName(name).
		Configure(
			// Limit default buffering for edge connections.
			envoy_listeners.ConnectionBufferLimit(DefaultConnectionBuffer),
			// Roughly balance incoming connections.
			envoy_listeners.EnableReusePort(true),
			// Always sniff for TLS.
			envoy_listeners.TLSInspector(),
		), limits
}
