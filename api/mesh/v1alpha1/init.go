package v1alpha1

import (
	"encoding/gob"
)

func init() {
	gob.Register(&TrafficRoute_LoadBalancer_RoundRobin_{})
	gob.Register(&ProxyTemplate_Modifications_Listener_{})
	gob.Register(&ProxyTemplate_Modifications_Cluster_{})
	gob.Register(&ProxyTemplate_Modifications_NetworkFilter_{})
	gob.Register(&ProxyTemplate_Modifications_HttpFilter_{})
	gob.Register(&ProxyTemplate_Modifications_VirtualHost_{})
}
