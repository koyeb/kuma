package clusters

import (
	envoy_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
)

type PanicModeConfigurer struct {
	IsOverriden bool
	Value       float64
}

var _ ClusterConfigurer = &PanicModeConfigurer{}

func (u *PanicModeConfigurer) Configure(c *envoy_cluster.Cluster) error {
	if !u.IsOverriden {
		return nil
	}

	if c.CommonLbConfig == nil {
		c.CommonLbConfig = &envoy_cluster.Cluster_CommonLbConfig{}
	}
	c.CommonLbConfig.HealthyPanicThreshold = &envoy_type.Percent{Value: float64(u.Value)}

	return nil
}
