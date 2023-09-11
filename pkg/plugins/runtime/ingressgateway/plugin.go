package ingressgateway

import (
	"errors"

	config_core "github.com/kumahq/kuma/pkg/config/core"
	"github.com/kumahq/kuma/pkg/core"
	core_plugins "github.com/kumahq/kuma/pkg/core/plugins"
	"github.com/kumahq/kuma/pkg/plugins/runtime/ingressgateway/metadata"
	"github.com/kumahq/kuma/pkg/xds/generator"
	generator_core "github.com/kumahq/kuma/pkg/xds/generator/core"
	generator_secrets "github.com/kumahq/kuma/pkg/xds/generator/secrets"
	"github.com/kumahq/kuma/pkg/xds/template"
)

func init() {
	core_plugins.Register(metadata.PluginName, &plugin{})
}

var log = core.Log.WithName("plugin").WithName("runtime").WithName("ingress-gateway")

type plugin struct{}

var _ core_plugins.BootstrapPlugin = &plugin{}

func (p *plugin) BeforeBootstrap(context *core_plugins.MutablePluginContext, config core_plugins.PluginConfig) error {
	if context.Config().Environment == config_core.KubernetesEnvironment {
		return errors.New("kubernetes is unsupported")
	}

	return nil
}

func (p *plugin) AfterBootstrap(context *core_plugins.MutablePluginContext, config core_plugins.PluginConfig) error {
	// NOTE(nicoche): not sure if this is useful
	// Insert our resolver before the default so that we can intercept
	// builtin gateway dataplanes.
	generator.DefaultTemplateResolver = template.SequentialResolver(
		TemplateResolver{},
		generator.DefaultTemplateResolver,
	)

	generator.RegisterProfile(metadata.ProfileIngressGatewayProxy, NewProxyProfile(context.Config().Multizone.Zone.Name))

	log.Info("registered ingress-gateway plugin")
	return nil
}

func (p *plugin) Name() core_plugins.PluginName {
	return metadata.PluginName
}

func (p *plugin) Order() int {
	// NOTE(nicoche): this looks like a k8s setting. Take the same thing as the gateway.
	// It has to go before Environment is prepared, so we have resources registered in K8S schema
	return core_plugins.EnvironmentPreparingOrder - 1
}

func NewGenerator(zone string) Generator {
	return Generator{
		Zone:                     zone,
		HTTPFilterChainGenerator: &HTTPFilterChainGenerator{},
		ClusterGenerator:         &ClusterGenerator{},
	}
}

// NewProxyProfile returns a new resource generator profile for
// ingress gateway dataplanes.
func NewProxyProfile(zone string) generator_core.ResourceGenerator {
	return generator_core.CompositeResourceGenerator{
		generator.AdminProxyGenerator{},
		generator.PrometheusEndpointGenerator{},
		generator.TracingProxyGenerator{},
		generator.TransparentProxyGenerator{},
		generator.DNSGenerator{},
		NewGenerator(zone),
		generator_secrets.Generator{},
	}
}
