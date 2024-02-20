package mesh

import (
	"fmt"
	"golang.org/x/exp/maps"
	"hash/fnv"
	"net"
	"sort"
	"strconv"

	"github.com/kumahq/kuma/pkg/core/resources/model"
)

func (r *ZoneIngressResource) UsesInboundInterface(address net.IP, port uint32) bool {
	if r == nil {
		return false
	}
	if port == r.Spec.GetNetworking().GetPort() && overlap(address, net.ParseIP(r.Spec.GetNetworking().GetAddress())) {
		return true
	}
	if port == r.Spec.GetNetworking().GetAdvertisedPort() && overlap(address, net.ParseIP(r.Spec.GetNetworking().GetAdvertisedAddress())) {
		return true
	}
	return false
}

func (r *ZoneIngressResource) IsRemoteIngress(localZone string) bool {
	if r.Spec.GetZone() == "" || r.Spec.GetZone() == localZone {
		return false
	}
	return true
}

func (r *ZoneIngressResource) HasPublicAddress() bool {
	if r == nil {
		return false
	}
	return r.Spec.GetNetworking().GetAdvertisedAddress() != "" && r.Spec.GetNetworking().GetAdvertisedPort() != 0
}

func (r *ZoneIngressResource) AdminAddress(defaultAdminPort uint32) string {
	if r == nil {
		return ""
	}
	ip := r.Spec.GetNetworking().GetAddress()
	adminPort := r.Spec.GetNetworking().GetAdmin().GetPort()
	if adminPort == 0 {
		adminPort = defaultAdminPort
	}
	return net.JoinHostPort(ip, strconv.FormatUint(uint64(adminPort), 10))
}

func (r *ZoneIngressResource) Hash() []byte {
	hasher := fnv.New128a()
	_, _ = hasher.Write(model.HashMeta(r))
	_, _ = hasher.Write([]byte(r.Spec.GetNetworking().GetAddress()))
	_, _ = hasher.Write([]byte(r.Spec.GetNetworking().GetAdvertisedAddress()))
	return hasher.Sum(nil)
}

func (r *ZoneIngressResource) HashForMesh(meshName string) []byte {
	hasher := fnv.New128a()

	// Hash meta but version
	meta := r.GetMeta()
	_, _ = hasher.Write([]byte(r.Descriptor().Name))
	_, _ = hasher.Write([]byte(meta.GetMesh()))
	_, _ = hasher.Write([]byte(meta.GetName()))

	// Hash spec
	_, _ = hasher.Write([]byte(r.Spec.GetNetworking().GetAddress()))
	_, _ = hasher.Write([]byte(r.Spec.GetNetworking().GetAdvertisedAddress()))

	// Hash changes to the mesh
	for _, service := range r.Spec.GetAvailableServices() {
		if service.GetMesh() != meshName {
			continue
		}

		_, _ = hasher.Write([]byte(service.GetMesh()))
		_, _ = hasher.Write([]byte(fmt.Sprintf("%d", service.GetInstances())))
		_, _ = hasher.Write([]byte(fmt.Sprintf("%t", service.GetExternalService())))

		tags := service.GetTags()
		tagsKeys := maps.Keys(tags)
		sort.Strings(tagsKeys)

		for _, key := range tagsKeys {
			_, _ = hasher.Write([]byte(key))
			_, _ = hasher.Write([]byte(tags[key]))
		}
	}

	return hasher.Sum(nil)
}
