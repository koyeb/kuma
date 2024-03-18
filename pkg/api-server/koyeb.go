package api_server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/koyeb/koyeb-api-client-go-internal/api/v1/koyeb"
	kuma_cp "github.com/kumahq/kuma/pkg/config/app/kuma-cp"
	config_core "github.com/kumahq/kuma/pkg/config/core"
	rest_errors "github.com/kumahq/kuma/pkg/core/rest/errors"
	koyebpkg "github.com/kumahq/kuma/pkg/koyeb"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func addIsPropagatedEndpoint(
	ws *restful.WebService,
	cfg *kuma_cp.Config,
) {
	if cfg.Mode == config_core.Global {
		return
	}

	var koyebApiClient *koyeb.APIClient
	if cfg.KoyebApiUrl == "" {
		log.Info("KoyebApiUrl is not set")
		return
	}

	koyebCfg := koyeb.NewConfiguration()
	koyebCfg.Servers[0].URL = cfg.KoyebApiUrl
	koyebApiClient = koyeb.NewAPIClient(koyebCfg)

	zone := cfg.Multizone.Zone.Name

	ws.Route(
		ws.GET("/koyeb/runtime").To(runtimeIngressGateway(koyebApiClient, zone)).Returns(200, "OK", nil),
	)
}

type KoyebRuntime struct {
	Zone string `json:"zone"`
	// key is hostname
	IngressGateways map[string]*IngressGatewayRuntime `json:"ingress_gateways"`
}

type IngressGatewayRuntime struct {
	GeneratedAt time.Time `json:"generated_at"`
}

type Prober struct {
	zone           string
	koyebClient    *koyeb.APIClient
	cachedAt       time.Time
	cachedResponse *KoyebRuntime
	igwIps         []string
}

func (p *Prober) refreshIgwIps(ctx context.Context) error {
	ips, err := p.listIgwIps(ctx)
	if err != nil {
		return errors.Wrap(err, "could not list igw IPs from koyeb catalog")
	}

	p.igwIps = ips
	return nil
}

func (p *Prober) listIgwIps(ctx context.Context) ([]string, error) {
	if os.Getenv("KOYEB_ENV") == "dev" { //TODO: remove this
		return []string{"127.0.0.1"}, nil
	}

	resp, _, err := p.koyebClient.CatalogDatacentersApi.ListDatacenters(ctx).Execute()
	if err != nil {
		return nil, err
	}

	IgwIps := []string{}
	for _, koyebDC := range *resp.Datacenters {
		if *koyebDC.Id == p.zone {
			ips, err := net.LookupIP(*koyebDC.Domain)
			if err != nil {
				return nil, errors.Wrapf(err, "could not resolve %s", *koyebDC.Domain)
			}

			for _, ip := range ips {
				IgwIps = append(IgwIps, ip.String())
			}
		}
	}

	return IgwIps, nil
}

func (p *Prober) Do(ctx context.Context) (*KoyebRuntime, error) {
	if p.cachedResponse != nil && time.Since(p.cachedAt) < 3*time.Second {
		return p.cachedResponse, nil
	}

	out := &KoyebRuntime{
		Zone:            p.zone,
		IngressGateways: map[string]*IngressGatewayRuntime{},
	}

	if len(p.igwIps) == 0 {
		return nil, errors.New("no igw to contact")
	}

	for _, ip := range p.igwIps {
		url := fmt.Sprintf("http://%s:%d%s", ip, koyebpkg.RuntimeInfoPort, koyebpkg.RuntimeInfoPath)
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrapf(err, "could not probe ingress gateway at url %s", url)
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, errors.Wrap(err, "could not read body")
		}

		ipRuntimeInfo := koyebpkg.RuntimeInfo{}
		err = json.Unmarshal([]byte(body), &ipRuntimeInfo)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal igw runtime response")
		}

		out.IngressGateways[ip] = &IngressGatewayRuntime{
			GeneratedAt: ipRuntimeInfo.GeneratedAt,
		}
	}

	p.cachedResponse = out
	p.cachedAt = time.Now()
	return out, nil
}

func runtimeIngressGateway(koyebClient *koyeb.APIClient, zone string) restful.RouteFunction {
	prober := &Prober{
		zone:        zone,
		koyebClient: koyebClient,
	}

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			err := prober.refreshIgwIps(context.Background())
			if err != nil {
				log.Error(err, "Could not refresh ingress gateways IPs")
			}
			select {
			case <-ticker.C:

			}
		}
	}()

	return func(request *restful.Request, response *restful.Response) {
		ctx := request.Request.Context()

		result, err := prober.Do(ctx)
		if err != nil {
			rest_errors.HandleError(request.Request.Context(), response, err, "Could not probe ingress gateways")
			return
		}

		err = response.WriteAsJson(result)
		if err != nil {
			rest_errors.HandleError(request.Request.Context(), response, err, "Could not write response")
			return
		}

		return
	}
}
