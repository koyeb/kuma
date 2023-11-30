package koyeb

import (
	"context"
	"net/http"

	"github.com/koyeb/koyeb-api-client-go-internal/api/v1/koyeb"
)

type staticDatacenterCatalog struct{}

func (c staticDatacenterCatalog) ListDatacenters(ctx context.Context) koyeb.ApiListDatacentersRequest {
	return koyeb.ApiListDatacentersRequest{}
}

func (c staticDatacenterCatalog) ListDatacentersExecute(r koyeb.ApiListDatacentersRequest) (koyeb.ListDatacentersReply, *http.Response, error) {
	rep := koyeb.ListDatacentersReply{
		Datacenters: &[]koyeb.DatacenterListItem{
			{
				Id:          strPointer("par1"),
				RegionId:    strPointer("par"),
				Domain:      strPointer("localhost"),
				Coordinates: &[]string{"2.3522", "48.8566"},
			},
			{
				Id:          strPointer("fra1"),
				RegionId:    strPointer("fra"),
				Domain:      strPointer("localhost"),
				Coordinates: &[]string{"8.6821", "50.1109"},
			},
		},
	}

	return rep, nil, nil
}

func NewLocalDatacenterCatalog() koyeb.CatalogDatacentersApi {
	return staticDatacenterCatalog{}
}

func NewDatacenterCatalog(apiURL string) koyeb.CatalogDatacentersApi {
	cfg := koyeb.NewConfiguration()
	cfg.Servers[0].URL = apiURL

	return koyeb.NewAPIClient(cfg).CatalogDatacentersApi
}

type staticInternalDeploymentsApi struct{}

func (c staticInternalDeploymentsApi) ListAllRoutes(ctx context.Context) koyeb.ApiListAllRoutesRequest {
	return koyeb.ApiListAllRoutesRequest{}
}

func (c staticInternalDeploymentsApi) ListAllRoutesExecute(r koyeb.ApiListAllRoutesRequest) (koyeb.ListAllRoutesReply, *http.Response, error) {
	rep := koyeb.ListAllRoutesReply{
		Routes: &[]koyeb.ListAllRoutesReplyRoute{
			{
				Path:            strPointer(""),
				Domain:          strPointer("grpc.koyeb.app"),
				DeploymentGroup: strPointer("prod"),
				ServiceId:       strPointer("dp"),
				Port:            int32Pointer(8004),
				Datacenters: &[]string{
					"par1",
				},
				// Technically, the service `dp` cannot be a member of two
				// different apps. We take this shortcut to simplify the
				// local setup.
				AppId: strPointer("app-grpc"),
			},
			{
				Path:            strPointer("/http"),
				Domain:          strPointer("http.local.koyeb.app"),
				DeploymentGroup: strPointer("prod"),
				ServiceId:       strPointer("dp"),
				Port:            int32Pointer(8001),
				Datacenters: &[]string{
					"par1",
				},
				AppId: strPointer("app-http"),
			},
			{
				Path:            strPointer("/http2"),
				Domain:          strPointer("http.local.koyeb.app"),
				DeploymentGroup: strPointer("prod"),
				ServiceId:       strPointer("dp"),
				Port:            int32Pointer(8002),
				Datacenters: &[]string{
					"par1",
				},
				AppId: strPointer("app-http"),
			},
			{
				Path:            strPointer("/grpc"),
				Domain:          strPointer("http.local.koyeb.app"),
				DeploymentGroup: strPointer("prod"),
				ServiceId:       strPointer("dp"),
				Port:            int32Pointer(8004),
				Datacenters: &[]string{
					"par1",
				},
				AppId: strPointer("app-http"),
			},
			{
				Path:            strPointer("/ws"),
				Domain:          strPointer("http.local.koyeb.app"),
				DeploymentGroup: strPointer("prod"),
				ServiceId:       strPointer("dp"),
				Port:            int32Pointer(8011),
				Datacenters: &[]string{
					"par1",
				},
				AppId: strPointer("app-http"),
			},
		},
	}

	return rep, nil, nil
}

func NewLocalInternalDeployments() koyeb.InternalDeploymentsApi {
	return staticInternalDeploymentsApi{}
}

func NewInternalDeployments(apiURL string) koyeb.InternalDeploymentsApi {
	cfg := koyeb.NewConfiguration()
	cfg.Servers[0].URL = apiURL

	return koyeb.NewAPIClient(cfg).InternalDeploymentsApi
}

func strPointer(s string) *string {
	return &s
}

func int32Pointer(i int32) *int32 {
	return &i
}
