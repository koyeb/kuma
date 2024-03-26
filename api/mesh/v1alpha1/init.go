package v1alpha1

import (
	"encoding/gob"
)

func init() {
	gob.Register(&TrafficRoute_LoadBalancer_RoundRobin_{})
	gob.Register(&ProxyTemplate_Modifications_Listener_{})
}
