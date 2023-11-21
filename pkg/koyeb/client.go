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

func strPointer(s string) *string {
	return &s
}

func NewLocalDatacenterCatalog() koyeb.CatalogDatacentersApi {
	return staticDatacenterCatalog{}
}

func NewDatacenterCatalog(apiURL string) koyeb.CatalogDatacentersApi {
	cfg := koyeb.NewConfiguration()
	cfg.Servers[0].URL = apiURL

	return koyeb.NewAPIClient(cfg).CatalogDatacentersApi
}
