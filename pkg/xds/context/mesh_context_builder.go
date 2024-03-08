package context

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core"
	"github.com/kumahq/kuma/pkg/core/datasource"
	"github.com/kumahq/kuma/pkg/core/dns/lookup"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/apis/system"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/registry"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/core/runtime/component"
	"github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/dns/vips"
	"github.com/kumahq/kuma/pkg/events"
	"github.com/kumahq/kuma/pkg/log"
	xds_topology "github.com/kumahq/kuma/pkg/xds/topology"
)

var logger = core.Log.WithName("xds").WithName("context")

type safeChangedTypes struct {
	sync.Mutex
	types map[core_model.ResourceType]struct{}
}

func (s *safeChangedTypes) markTypeChanged(t core_model.ResourceType) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.types[t] = struct{}{}
}

func (s *safeChangedTypes) hasChanged(t core_model.ResourceType) bool {
	if s == nil {
		return true
	}

	s.Lock()
	defer s.Unlock()

	_, ok := s.types[t]
	return !ok
}

func (s *safeChangedTypes) unsafeClear(t core_model.ResourceType) {
	delete(s.types, t)
}

type safeChangedTypesByMesh struct {
	sync.Mutex
	v map[string]map[core_model.ResourceType]struct{}
}

func (s *safeChangedTypesByMesh) forMesh(mesh string) map[core_model.ResourceType]struct{} {
	if s == nil {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	got, ok := s.v[mesh]
	if !ok {
		return nil
	}

	return got
}

func (s *safeChangedTypesByMesh) clearMesh(mesh string) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.v[mesh] = map[core_model.ResourceType]struct{}{}
}

func (s *safeChangedTypesByMesh) markTypeChanged(mesh string, t core_model.ResourceType) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	_, ok := s.v[mesh]
	if !ok {
		s.v[mesh] = map[core_model.ResourceType]struct{}{}
	}

	s.v[mesh][t] = struct{}{}
}

type meshContextBuilder struct {
	rm                       manager.ReadOnlyResourceManager
	typeSet                  map[core_model.ResourceType]struct{}
	ipFunc                   lookup.LookupIPFunc
	zone                     string
	vipsPersistence          *vips.Persistence
	topLevelDomain           string
	vipPort                  uint32
	rsGraphBuilder           ReachableServicesGraphBuilder
	globalChangedTypes       *safeChangedTypes
	safeChangedTypesByMesh   *safeChangedTypesByMesh
	eventBus                 events.EventBus
	hashCacheBaseMeshContext *cache.Cache
	hashCacheGlobalContext   *cache.Cache
}

// MeshContextBuilder
type MeshContextBuilder interface {
	Build(ctx context.Context, meshName string) (MeshContext, error)

	// BuildGlobalContextIfChanged builds GlobalContext only if `latest` is nil or hash is different
	// If hash is the same, the return `latest`
	BuildGlobalContextIfChanged(ctx context.Context, latest *GlobalContext, meshName string) (*GlobalContext, error)

	// BuildBaseMeshContextIfChanged builds BaseMeshContext only if `latest` is nil or hash is different
	// If hash is the same, the return `latest`
	BuildBaseMeshContextIfChanged(ctx context.Context, meshName string, latest *BaseMeshContext) (*BaseMeshContext, error)
	BuildBaseMeshContextIfChangedV2(ctx context.Context, meshName string, latest *BaseMeshContext) (*BaseMeshContext, error)

	// BuildIfChanged builds MeshContext only if latestMeshCtx is nil or hash of
	// latestMeshCtx is different.
	// If hash is the same, then the function returns the passed latestMeshCtx.
	// Hash returned in MeshContext can never be empty.
	BuildIfChanged(ctx context.Context, meshName string, latestMeshCtx *MeshContext) (*MeshContext, error)

	Start(stop <-chan struct{}) error
	NeedLeaderElection() bool
}

// meshContextCleanupTime is the time after which the mesh context is removed from
// the longer TTL cache.
// It exists to ensure contexts of deleted Meshes are eventually cleaned up.
const meshContextCleanupTime = 45 * time.Minute
const globalContextCleanupTime = 2 * time.Minute

type MeshContextBuilderComponent interface {
	MeshContextBuilder
	component.Component
}

func NewMeshContextBuilderComponent(
	rm manager.ReadOnlyResourceManager,
	types []core_model.ResourceType, // types that should be taken into account when MeshContext is built.
	ipFunc lookup.LookupIPFunc,
	zone string,
	vipsPersistence *vips.Persistence,
	topLevelDomain string,
	vipPort uint32,
	rsGraphBuilder ReachableServicesGraphBuilder,
	eventBus events.EventBus,
) MeshContextBuilderComponent {
	typeSet := map[core_model.ResourceType]struct{}{}
	for _, typ := range types {
		typeSet[typ] = struct{}{}
	}

	return &meshContextBuilder{
		rm:              rm,
		typeSet:         typeSet,
		ipFunc:          ipFunc,
		zone:            zone,
		vipsPersistence: vipsPersistence,
		topLevelDomain:  topLevelDomain,
		vipPort:         vipPort,
		rsGraphBuilder:  rsGraphBuilder,
		safeChangedTypesByMesh: &safeChangedTypesByMesh{
			v: map[string]map[core_model.ResourceType]struct{}{},
		},
		globalChangedTypes: &safeChangedTypes{
			types: map[core_model.ResourceType]struct{}{},
		},
		eventBus:                 eventBus,
		hashCacheBaseMeshContext: cache.New(meshContextCleanupTime, time.Duration(int64(float64(meshContextCleanupTime)*0.9))),
		hashCacheGlobalContext:   cache.New(globalContextCleanupTime, time.Duration(int64(float64(globalContextCleanupTime)*0.9))),
	}
}

func NewMeshContextBuilder(
	rm manager.ReadOnlyResourceManager,
	types []core_model.ResourceType, // types that should be taken into account when MeshContext is built.
	ipFunc lookup.LookupIPFunc,
	zone string,
	vipsPersistence *vips.Persistence,
	topLevelDomain string,
	vipPort uint32,
	rsGraphBuilder ReachableServicesGraphBuilder,
) MeshContextBuilder {
	typeSet := map[core_model.ResourceType]struct{}{}
	for _, typ := range types {
		typeSet[typ] = struct{}{}
	}

	return &meshContextBuilder{
		rm:                       rm,
		typeSet:                  typeSet,
		ipFunc:                   ipFunc,
		zone:                     zone,
		vipsPersistence:          vipsPersistence,
		topLevelDomain:           topLevelDomain,
		vipPort:                  vipPort,
		rsGraphBuilder:           rsGraphBuilder,
		safeChangedTypesByMesh:   nil,
		hashCacheBaseMeshContext: nil,
		hashCacheGlobalContext:   nil,
	}
}

func useReactiveBuildBaseMeshContext() bool {
	return os.Getenv("EXPERIMENTAL_REACTIVE_BASE_MESH_CONTEXT") != ""
}

func (m *meshContextBuilder) Start(stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	l := log.AddFieldsFromCtx(logger, ctx, context.Background())

	listener := m.eventBus.Subscribe(func(event events.Event) bool {
		resChange, ok := event.(events.ResourceChangedEvent)
		if !ok {
			return false
		}

		// if resChange.TenantID != tenantID {
		// 	return false
		// }

		_, ok = m.typeSet[resChange.Type]
		return ok
	})

	for {
		select {
		case <-stop:
			return nil

		case event := <-listener.Recv():

			if useReactiveBuildBaseMeshContext() {
				resChange := event.(events.ResourceChangedEvent)
				l.Info("Received", "ResourceChangedEvent", resChange)
				mesh := resChange.Key.Mesh

				if mesh != "" {
					l.Info("Type has changed for mesh", "type", resChange.Type, "mesh", mesh)
					m.setTypeChanged(mesh, resChange.Type)
				} else {
					desc, err := registry.Global().DescriptorFor(resChange.Type)
					if err != nil {
						l.Error(err, "Could not get type")
						continue
					}
					if desc.Name == system.ConfigType {
						mesh, ok := vips.MeshFromConfigKey(resChange.Key.Name)
						if !ok {
							continue
						}

						l.Info("Type has changed for mesh", "type", "Config", "mesh", mesh)
						m.setTypeChanged(mesh, resChange.Type)
						// well should alwys be true because we're in a branch where mesh == ""?
					} else if desc.Scope == core_model.ScopeGlobal {
						m.globalChangedTypes.markTypeChanged(resChange.Type)
					}

				}

			}
		}
	}
}

func (m *meshContextBuilder) NeedLeaderElection() bool {
	return false
}

func (m *meshContextBuilder) Build(ctx context.Context, meshName string) (MeshContext, error) {
	meshCtx, err := m.BuildIfChanged(ctx, meshName, nil)
	if err != nil {
		return MeshContext{}, err
	}
	return *meshCtx, nil
}

func (m *meshContextBuilder) shouldLogExcessively(meshName string) bool {
	for _, mesh := range strings.Split(os.Getenv("DEBUG_MESHES"), ",") {
		if mesh == meshName {
			return true
		}
	}

	return false
}

func (m *meshContextBuilder) BuildIfChanged(ctx context.Context, meshName string, latestMeshCtx *MeshContext) (*MeshContext, error) {
	l := log.AddFieldsFromCtx(logger, ctx, context.Background())
	isDefaultMesh := meshName == "default"

	if m.shouldLogExcessively(meshName) {
		l.Info("Running BuildIfChanged", "mesh", meshName)
	}

	var latestGlobalContext *GlobalContext
	cachedGlobalContext, ok := m.hashCacheGlobalContext.Get("")
	if ok && !isDefaultMesh {
		latestGlobalContext = cachedGlobalContext.(*GlobalContext)
	}
	globalContext, err := m.BuildGlobalContextIfChanged(ctx, latestGlobalContext, meshName)
	if err != nil {
		return nil, err
	}
	m.hashCacheGlobalContext.SetDefault("", globalContext)

	var baseMeshContext *BaseMeshContext
	if useReactiveBuildBaseMeshContext() && !isDefaultMesh {

		// Check hashCache first for an existing mesh latestContext
		var latestBaseMeshContext *BaseMeshContext
		if m.hashCacheBaseMeshContext != nil {
			if cached, ok := m.hashCacheBaseMeshContext.Get(meshName); ok {
				latestBaseMeshContext = cached.(*BaseMeshContext)
				if m.shouldLogExcessively(meshName) {
					l.Info("Found latest base mesh context to re-use", "mesh", meshName, "hash", latestBaseMeshContext.hash)
				}
			}
		}

		baseMeshContext, err = m.BuildBaseMeshContextIfChangedV2(ctx, meshName, latestBaseMeshContext)
		if err != nil {
			return nil, err
		}
	} else {
		baseMeshContext, err = m.BuildBaseMeshContextIfChanged(ctx, meshName, nil)
		if err != nil {
			return nil, err
		}
	}

	// By always setting the mesh context, we refresh the TTL
	// with the effect that often used contexts remain in the cache while no
	// longer used contexts are evicted.
	if m.shouldLogExcessively(meshName) {
		l.Info("Saving base mesh context in hash cache", "mesh", meshName, "hash", baseMeshContext.hash)
	}
	m.hashCacheBaseMeshContext.SetDefault(meshName, baseMeshContext)

	var managedTypes []core_model.ResourceType // The types not managed by global nor baseMeshContext
	resources := NewResources()
	// Build all the local entities from the parent contexts
	for resType := range m.typeSet {
		rl, ok := globalContext.ResourceMap[resType]
		if ok { // Exists in global context take it from there
			switch resType {
			case core_mesh.MeshType: // Remove our own mesh from the list
				otherMeshes := rl.NewItem().Descriptor().NewList()
				for _, rentry := range rl.GetItems() {
					if rentry.GetMeta().GetName() != meshName {
						err := otherMeshes.AddItem(rentry)
						if err != nil {
							return nil, err
						}
					}
				}
				otherMeshes.GetPagination().SetTotal(uint32(len(otherMeshes.GetItems())))
				rl = otherMeshes
			}
			resources.MeshLocalResources[resType] = rl
		} else {
			rl, ok = baseMeshContext.ResourceMap[resType]
			if ok { // Exist in the baseMeshContext take it from there
				resources.MeshLocalResources[resType] = rl
			} else { // absent from all parent contexts get it now
				managedTypes = append(managedTypes, resType)
				rl, err = m.fetchResourceList(ctx, resType, baseMeshContext.Mesh, nil)
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("could not fetch resources of type:%s", resType))
				}
				resources.MeshLocalResources[resType] = rl
			}
		}
	}

	if isDefaultMesh {
		if err := m.decorateWithCrossMeshResources(ctx, resources); err != nil {
			return nil, errors.Wrap(err, "failed to retrieve cross mesh resources")
		}
	}

	// This base64 encoding seems superfluous but keeping it for backward compatibility
	newHash := base64.StdEncoding.EncodeToString(m.hash(meshName, globalContext, baseMeshContext, managedTypes, resources))
	if m.shouldLogExcessively(meshName) {
		l.Info("Computed mesh context hash", "mesh", meshName, "hash", newHash)
	}
	if !isDefaultMesh && latestMeshCtx != nil && newHash == latestMeshCtx.Hash {
		if m.shouldLogExcessively(meshName) {
			l.Info("Latest mesh context hash is the same as the hash of resources needed to compute the current mesh context. Returning it now", "mesh", meshName, "hash", latestMeshCtx.Hash)
		}
		return latestMeshCtx, nil
	}

	if latestMeshCtx == nil {
		l.Info("mesh context never computed for this mesh, computing", "mesh", meshName)
	} else {
		l.Info("latest mesh context hash is different than computed hash, recomputing", "mesh", meshName)
	}

	dataplanes := resources.Dataplanes().Items
	dataplanesByName := make(map[string]*core_mesh.DataplaneResource, len(dataplanes))
	for _, dp := range dataplanes {
		dataplanesByName[dp.Meta.GetName()] = dp
	}

	virtualOutboundView, err := m.vipsPersistence.GetByMesh(ctx, meshName)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch vips")
	}
	// resolve all the domains
	domains, outbounds := xds_topology.VIPOutbounds(virtualOutboundView, m.topLevelDomain, m.vipPort)

	mesh := baseMeshContext.Mesh
	zoneIngresses := resources.ZoneIngresses().Items
	zoneEgresses := resources.ZoneEgresses().Items
	externalServices := resources.ExternalServices().Items
	endpointMap := xds_topology.BuildEdsEndpointMap(mesh, m.zone, dataplanes, zoneIngresses, zoneEgresses, externalServices)

	crossMeshEndpointMap := map[string]xds.EndpointMap{}
	for otherMeshName, gateways := range resources.gatewaysAndDataplanesForMesh(mesh) {
		crossMeshEndpointMap[otherMeshName] = xds_topology.BuildCrossMeshEndpointMap(
			mesh,
			gateways.Mesh,
			m.zone,
			gateways.Gateways,
			gateways.Dataplanes,
			zoneIngresses,
			zoneEgresses,
		)
	}

	return &MeshContext{
		Hash:                   newHash,
		Resource:               mesh,
		Resources:              resources,
		DataplanesByName:       dataplanesByName,
		EndpointMap:            endpointMap,
		CrossMeshEndpoints:     crossMeshEndpointMap,
		VIPDomains:             domains,
		VIPOutbounds:           outbounds,
		ServiceTLSReadiness:    m.resolveTLSReadiness(mesh, resources.ServiceInsights()),
		DataSourceLoader:       datasource.NewStaticLoader(resources.Secrets().Items),
		ReachableServicesGraph: m.rsGraphBuilder(meshName, resources),
	}, nil
}

func (m *meshContextBuilder) BuildGlobalContextIfChanged(ctx context.Context, latest *GlobalContext, meshName string) (*GlobalContext, error) {
	rmap := ResourceMap{}
	for t := range m.typeSet {
		desc, err := registry.Global().DescriptorFor(t)
		if err != nil {
			return nil, err
		}

		// Only pick the global stuff
		if desc.Scope != core_model.ScopeGlobal {
			continue
		}

		// For config we ignore them atm and prefer to rely on more specific filters.
		if desc.Name == system.ConfigType {
			continue
		}

		// If no change has been registered for this type, pick the resource from the latest global context,
		// else fetch it from store
		if !m.globalChangedTypes.hasChanged(t) && latest != nil {
			rmap[t] = latest.ResourceMap[t]
		} else {
			m.globalChangedTypes.Lock()
			rmap[t], err = m.fetchResourceList(ctx, t, nil, nil)
			if err != nil {
				m.globalChangedTypes.Unlock()
				return nil, errors.Wrap(err, "failed to build global context")
			}
			m.globalChangedTypes.unsafeClear(t)
			m.globalChangedTypes.Unlock()
		}
	}

	newHash := rmap.HashForMesh(meshName)
	if meshName != "default" && latest != nil && bytes.Equal(newHash, latest.hash) {
		return latest, nil
	}
	return &GlobalContext{
		hash:        newHash,
		ResourceMap: rmap,
	}, nil
}

func (m *meshContextBuilder) BuildBaseMeshContextIfChangedV2(ctx context.Context, meshName string, latest *BaseMeshContext) (*BaseMeshContext, error) {
	l := log.AddFieldsFromCtx(logger, ctx, context.Background())
	if latest == nil {
		l.Info("no latest base mesh context to use or not using reactive method. Fallback to default BuildBaseMeshContextIfChanged", "mesh", meshName)
		m.clearTypeChanged(meshName)
		return m.BuildBaseMeshContextIfChanged(ctx, meshName, nil)
	}

	changedTypes := m.safeChangedTypesByMesh.forMesh(meshName)
	if changedTypes == nil || len(changedTypes) == 0 {
		// No occurence of this mesh in changed types. Let's re-use latest base mesh context
		// l.Info("no resource changed, re-using latest base mesh context to build mesh context", "mesh", meshName)
		if m.shouldLogExcessively(meshName) {
			l.Info("No changed types for this mesh. Re-using latest base mesh context", "mesh", meshName)
		}
		return latest, nil
	}

	rmap := ResourceMap{}

	// Find mesh, either in last base context if it hasn't changed since or in the store
	mesh := core_mesh.NewMeshResource()
	_, meshChanged := changedTypes[core_mesh.MeshType]
	if !meshChanged && latest != nil {
		meshList := latest.ResourceMap[core_mesh.MeshType].(*core_mesh.MeshResourceList)
		mesh = meshList.Items[0]
		if m.shouldLogExcessively(meshName) {
			l.Info("Found mesh in latest base mesh context's resource map", "mesh", meshName)
		}
	} else {
		if err := m.rm.Get(ctx, mesh, core_store.GetByKey(meshName, core_model.NoMesh)); err != nil {
			return nil, errors.Wrapf(err, "could not fetch mesh %s", meshName)
		}
		if m.shouldLogExcessively(meshName) {
			l.Info("Found mesh in the store", "mesh", meshName)
		}
	}

	// Add mesh to resource map
	rmap[core_mesh.MeshType] = mesh.Descriptor().NewList()
	_ = rmap[core_mesh.MeshType].AddItem(mesh)

	for t := range m.typeSet {
		desc, err := registry.Global().DescriptorFor(t)
		if err != nil {
			return nil, err
		}

		switch {
		case desc.IsPolicy || desc.Name == core_mesh.MeshGatewayType || desc.Name == core_mesh.ExternalServiceType || desc.Name == system.SecretType || desc.Name == core_mesh.DataplaneType:
			rmap[t], err = m.fetchResourceListIfChanged(ctx, latest, t, mesh, nil)
		case desc.Name == system.ConfigType:
			rmap[t], err = m.fetchResourceListIfChanged(ctx, latest, t, mesh, func(rs core_model.Resource) bool {
				return rs.GetMeta().GetName() == vips.ConfigKey(meshName)
			})
		default:
			// Do nothing we're not interested in this type
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to build base mesh context")
		}
	}

	// Reset changed types for this mesh
	if m.shouldLogExcessively(meshName) {
		l.Info("Clear types changed", "mesh", meshName)
	}
	m.clearTypeChanged(meshName)

	newHash := rmap.HashForMesh(meshName)
	if latest != nil && bytes.Equal(newHash, latest.hash) {
		return latest, nil
	}

	return &BaseMeshContext{
		hash:        newHash,
		Mesh:        mesh,
		ResourceMap: rmap,
	}, nil
}

func (m *meshContextBuilder) setTypeChanged(mesh string, t core_model.ResourceType) {
	m.safeChangedTypesByMesh.markTypeChanged(mesh, t)
}

func (m *meshContextBuilder) clearTypeChanged(mesh string) {
	m.safeChangedTypesByMesh.clearMesh(mesh)
}

func (m *meshContextBuilder) BuildBaseMeshContextIfChanged(ctx context.Context, meshName string, latest *BaseMeshContext) (*BaseMeshContext, error) {
	mesh := core_mesh.NewMeshResource()
	if err := m.rm.Get(ctx, mesh, core_store.GetByKey(meshName, core_model.NoMesh)); err != nil {
		return nil, errors.Wrapf(err, "could not fetch mesh %s", meshName)
	}
	rmap := ResourceMap{}
	// Add the mesh to the resourceMap
	rmap[core_mesh.MeshType] = mesh.Descriptor().NewList()
	_ = rmap[core_mesh.MeshType].AddItem(mesh)
	rmap[core_mesh.MeshType].GetPagination().SetTotal(1)
	for t := range m.typeSet {
		desc, err := registry.Global().DescriptorFor(t)
		if err != nil {
			return nil, err
		}

		switch {
		case desc.IsPolicy || desc.Name == core_mesh.MeshGatewayType || desc.Name == core_mesh.ExternalServiceType || desc.Name == system.SecretType || desc.Name == core_mesh.DataplaneType:
			rmap[t], err = m.fetchResourceList(ctx, t, mesh, nil)
		case desc.Name == system.ConfigType:
			rmap[t], err = m.fetchResourceList(ctx, t, mesh, func(rs core_model.Resource) bool {
				return rs.GetMeta().GetName() == vips.ConfigKey(meshName)
			})
		default:
			// DO nothing we're not interested in this type
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to build base mesh context")
		}
	}
	newHash := rmap.HashForMesh(meshName)
	if latest != nil && bytes.Equal(newHash, latest.hash) {
		return latest, nil
	}
	return &BaseMeshContext{
		hash:        newHash,
		Mesh:        mesh,
		ResourceMap: rmap,
	}, nil
}

type filterFn = func(rs core_model.Resource) bool

// fetch resource from latest base mesh context if it hasn't changed. Else, pull it from the store
func (m *meshContextBuilder) fetchResourceListIfChanged(ctx context.Context, latest *BaseMeshContext, resType core_model.ResourceType, mesh *core_mesh.MeshResource, filterFn filterFn) (core_model.ResourceList, error) {
	l := log.AddFieldsFromCtx(logger, ctx, context.Background())
	meshName := mesh.GetMeta().GetName()

	if latest == nil {
		if m.shouldLogExcessively(meshName) {
			l.Info("Fetching from store", "mesh", meshName, "resourcetype", resType)
		}
		return m.fetchResourceList(ctx, resType, mesh, filterFn)
	}

	changedTypes := m.safeChangedTypesByMesh.forMesh(meshName)
	if changedTypes == nil || len(changedTypes) == 0 {
		if m.shouldLogExcessively(meshName) {
			l.Info("Resource not found in changedTypes list. Using version from latest base mesh context", "mesh", meshName, "resourcetype", resType)
		}
		return latest.ResourceMap[core_mesh.MeshType], nil

		_, hasChanged := changedTypes[resType]
		if !hasChanged {
			if m.shouldLogExcessively(meshName) {
				l.Info("Resource has not changed since last time. Using version from latest base mesh context", "mesh", meshName, "resourcetype", resType)
			}
			return latest.ResourceMap[resType], nil
		}

		if m.shouldLogExcessively(meshName) {
			l.Info("Resource has not changed since last time. Using version from latest base mesh context", "mesh", meshName, "resourcetype", resType)
		}
		return latest.ResourceMap[resType], nil
	}

	if m.shouldLogExcessively(meshName) {
		l.Info("Fetching from store", "mesh", meshName, "resourcetype", resType)
	}
	return m.fetchResourceList(ctx, resType, mesh, filterFn)
}

// fetch all resources of a type with potential filters etc
func (m *meshContextBuilder) fetchResourceList(ctx context.Context, resType core_model.ResourceType, mesh *core_mesh.MeshResource, filterFn filterFn) (core_model.ResourceList, error) {
	l := log.AddFieldsFromCtx(logger, ctx, context.Background())
	var listOptsFunc []core_store.ListOptionsFunc
	desc, err := registry.Global().DescriptorFor(resType)
	if err != nil {
		return nil, err
	}
	switch desc.Scope {
	case core_model.ScopeGlobal:
	case core_model.ScopeMesh:
		if mesh != nil {
			listOptsFunc = append(listOptsFunc, core_store.ListByMesh(mesh.GetMeta().GetName()))
		}
	default:
		return nil, fmt.Errorf("unknown resource scope:%s", desc.Scope)
	}

	// For some resources we apply extra filters
	switch resType {
	case core_mesh.ServiceInsightType:
		if mesh == nil {
			return desc.NewList(), nil
		}
		// ServiceInsights in XDS generation are only used to check whether the destination is ready to receive mTLS traffic.
		// This information is only useful when mTLS is enabled with PERMISSIVE mode.
		// Not including this into mesh hash for other cases saves us unnecessary XDS config generations.
		if backend := mesh.GetEnabledCertificateAuthorityBackend(); backend == nil || backend.Mode == mesh_proto.CertificateAuthorityBackend_STRICT {
			return desc.NewList(), nil
		}
	}
	listOptsFunc = append(listOptsFunc, core_store.ListOrdered())
	list := desc.NewList()
	if err := m.rm.List(ctx, list, listOptsFunc...); err != nil {
		return nil, err
	}
	if resType != core_mesh.ZoneIngressType && resType != core_mesh.DataplaneType && filterFn == nil {
		// No post processing stuff so return the list as is
		return list, nil
	}
	list, err = modifyAllEntries(list, func(resource core_model.Resource) (core_model.Resource, error) {
		// Because we're not using the pagination store we need to do the filtering ourselves outside of the store
		// I believe this is to maximize cachability of the store
		if filterFn != nil && !filterFn(resource) {
			return nil, nil
		}
		switch resType {
		case core_mesh.ZoneIngressType:
			zi, ok := resource.(*core_mesh.ZoneIngressResource)
			if !ok {
				return nil, errors.New("entry is not a zoneIngress this shouldn't happen")
			}
			zi, err := xds_topology.ResolveZoneIngressPublicAddress(m.ipFunc, zi)
			if err != nil {
				l.Error(err, "failed to resolve zoneIngress's domain name, ignoring zoneIngress", "name", zi.GetMeta().GetName())
				return nil, nil
			}
			return zi, nil
		case core_mesh.DataplaneType:
			list, err = modifyAllEntries(list, func(resource core_model.Resource) (core_model.Resource, error) {
				dp, ok := resource.(*core_mesh.DataplaneResource)
				if !ok {
					return nil, errors.New("entry is not a dataplane this shouldn't happen")
				}
				zi, err := xds_topology.ResolveDataplaneAddress(m.ipFunc, dp)
				if err != nil {
					l.Error(err, "failed to resolve dataplane's domain name, ignoring dataplane", "mesh", dp.GetMeta().GetMesh(), "name", dp.GetMeta().GetName())
					return nil, nil
				}
				return zi, nil
			})
		}
		return resource, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// takes a resourceList and modify it as needed
func modifyAllEntries(list core_model.ResourceList, fn func(resource core_model.Resource) (core_model.Resource, error)) (core_model.ResourceList, error) {
	newList := list.NewItem().Descriptor().NewList()
	for _, v := range list.GetItems() {
		ni, err := fn(v)
		if err != nil {
			return nil, err
		}
		if ni != nil {
			err := newList.AddItem(ni)
			if err != nil {
				return nil, err
			}
		}
	}
	newList.GetPagination().SetTotal(uint32(len(newList.GetItems())))
	return newList, nil
}

func (m *meshContextBuilder) resolveTLSReadiness(mesh *core_mesh.MeshResource, serviceInsights *core_mesh.ServiceInsightResourceList) map[string]bool {
	tlsReady := map[string]bool{}

	backend := mesh.GetEnabledCertificateAuthorityBackend()
	// TLS readiness is irrelevant unless we are using PERMISSIVE TLS, so skip
	// checking ServiceInsights if we aren't.
	if backend == nil || backend.Mode != mesh_proto.CertificateAuthorityBackend_PERMISSIVE {
		return tlsReady
	}

	if len(serviceInsights.Items) == 0 {
		// Nothing about the TLS readiness has been reported yet
		logger.Info("could not determine service TLS readiness, ServiceInsight is not yet present")
		return tlsReady
	}

	for svc, insight := range serviceInsights.Items[0].Spec.GetServices() {
		if insight.ServiceType == mesh_proto.ServiceInsight_Service_external {
			tlsReady[svc] = true
		} else {
			tlsReady[svc] = insight.IssuedBackends[backend.Name] == (insight.Dataplanes.Offline + insight.Dataplanes.Online)
		}
	}
	return tlsReady
}

func (m *meshContextBuilder) decorateWithCrossMeshResources(ctx context.Context, resources Resources) error {
	// Expand with crossMesh info
	otherMeshesByName := map[string]*core_mesh.MeshResource{}
	for _, m := range resources.OtherMeshes().GetItems() {
		otherMeshesByName[m.GetMeta().GetName()] = m.(*core_mesh.MeshResource)
		resources.CrossMeshResources[m.GetMeta().GetName()] = map[core_model.ResourceType]core_model.ResourceList{
			core_mesh.DataplaneType:   &core_mesh.DataplaneResourceList{},
			core_mesh.MeshGatewayType: &core_mesh.MeshGatewayResourceList{},
		}
	}
	var gatewaysByMesh map[string]core_model.ResourceList
	if _, ok := m.typeSet[core_mesh.MeshGatewayType]; ok {
		// For all meshes, get all cross mesh gateways
		rl, err := m.fetchResourceList(ctx, core_mesh.MeshGatewayType, nil, func(rs core_model.Resource) bool {
			_, exists := otherMeshesByName[rs.GetMeta().GetMesh()]
			return exists && rs.(*core_mesh.MeshGatewayResource).Spec.IsCrossMesh()
		})
		if err != nil {
			return errors.Wrap(err, "could not fetch cross mesh meshGateway resources")
		}
		gatewaysByMesh, err = core_model.ResourceListByMesh(rl)
		if err != nil {
			return errors.Wrap(err, "failed building cross mesh meshGateway resources")
		}
	}
	if _, ok := m.typeSet[core_mesh.DataplaneType]; ok {
		for otherMeshName, gws := range gatewaysByMesh { // Only iterate over meshes that have crossMesh gateways
			otherMesh := otherMeshesByName[otherMeshName]
			var gwResources []*core_mesh.MeshGatewayResource
			for _, gw := range gws.GetItems() {
				gwResources = append(gwResources, gw.(*core_mesh.MeshGatewayResource))
			}
			rl, err := m.fetchResourceList(ctx, core_mesh.DataplaneType, otherMesh, func(rs core_model.Resource) bool {
				dp := rs.(*core_mesh.DataplaneResource)
				if !dp.Spec.IsBuiltinGateway() {
					return false
				}
				return xds_topology.SelectGateway(gwResources, dp.Spec.Matches) != nil
			})
			if err != nil {
				return errors.Wrap(err, "could not fetch cross mesh meshGateway resources")
			}
			resources.CrossMeshResources[otherMeshName][core_mesh.DataplaneType] = rl
			resources.CrossMeshResources[otherMeshName][core_mesh.MeshGatewayType] = gws
		}
	}
	return nil
}

func (m *meshContextBuilder) hash(meshName string, globalContext *GlobalContext, baseMeshContext *BaseMeshContext, managedTypes []core_model.ResourceType, resources Resources) []byte {
	slices.Sort(managedTypes)
	hasher := fnv.New128a()
	_, _ = hasher.Write(globalContext.hash)
	_, _ = hasher.Write(baseMeshContext.hash)
	for _, resType := range managedTypes {
		_, _ = hasher.Write(core_model.ResourceListHashForMesh(resources.MeshLocalResources[resType], meshName))
	}

	return hasher.Sum(nil)
}
