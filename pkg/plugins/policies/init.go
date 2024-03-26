package policies

import (
	"encoding/gob"

	donothing_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/donothingpolicy/api/v1alpha1"
	accesslog_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshaccesslog/api/v1alpha1"
	circuitbreaker_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshcircuitbreaker/api/v1alpha1"
	meshfaultinjection_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshfaultinjection/api/v1alpha1"
	hc_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshhealthcheck/api/v1alpha1"
	httproute_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshhttproute/api/v1alpha1"
	meshloadbalancingstrategy_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshloadbalancingstrategy/api/v1alpha1"
	proxypatch_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshproxypatch/api/v1alpha1"
	ratelimit_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshratelimit/api/v1alpha1"
	retry_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshretry/api/v1alpha1"
	meshtcproute_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshtcproute/api/v1alpha1"
	timeout_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshtimeout/api/v1alpha1"
	trace_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshtrace/api/v1alpha1"
	meshtrafficpermission_v1alpha1 "github.com/kumahq/kuma/pkg/plugins/policies/meshtrafficpermission/api/v1alpha1"
)

func init() {
	gob.Register(&circuitbreaker_v1alpha1.MeshCircuitBreakerResourceList{})
	gob.Register(&meshtcproute_v1alpha1.MeshTCPRouteResourceList{})
	gob.Register(&meshtrafficpermission_v1alpha1.MeshTrafficPermissionResourceList{})
	gob.Register(&meshfaultinjection_v1alpha1.MeshFaultInjectionResourceList{})
	gob.Register(&meshloadbalancingstrategy_v1alpha1.MeshLoadBalancingStrategyResourceList{})
	gob.Register(&donothing_v1alpha1.DoNothingPolicyResourceList{})
	gob.Register(&hc_v1alpha1.MeshHealthCheckResourceList{})
	gob.Register(&accesslog_v1alpha1.MeshAccessLogResourceList{})
	gob.Register(&httproute_v1alpha1.MeshHTTPRouteResourceList{})
	gob.Register(&ratelimit_v1alpha1.MeshRateLimitResourceList{})
	gob.Register(&circuitbreaker_v1alpha1.MeshCircuitBreakerResourceList{})
	gob.Register(&proxypatch_v1alpha1.MeshProxyPatchResourceList{})
	gob.Register(&timeout_v1alpha1.MeshTimeoutResourceList{})
	gob.Register(&trace_v1alpha1.MeshTraceResourceList{})
	gob.Register(&retry_v1alpha1.MeshRetryResourceList{})
}
