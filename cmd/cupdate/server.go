package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/web"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ConfigureServer(
	config *Config,
	httpClient *httputil.Client,
	readStore *store.Store,
	worker *worker.Worker,
	processQueue *worker.Queue[oci.Reference],
	targetPlatform platform.ContinuousGrapher,
) *http.Server {
	logoProxy := api.CompoundProxy{
		Proxies: []api.LogoProxy{
			&api.LogoFSProxy{
				FS: os.DirFS(config.Logos.Path),
			},
			&api.LogoHTTPProxy{
				Client: httpClient,
				GetURL: readStore.GetImageLogo,
			},
		},
	}

	mux := http.NewServeMux()

	apiServer := api.NewServer(readStore, worker.Hub, processQueue, logoProxy, targetPlatform)
	apiServer.WebAddress = config.Web.Address
	mux.Handle("/api/v1/", apiServer)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/livez", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Figure out what checks to have
		w.WriteHeader(http.StatusOK)
	}))

	if !config.Web.Disabled {
		mux.Handle("/", web.MustNewEmbeddedServer())
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.API.Address, config.API.Port),
		Handler: http.NewCrossOriginProtection().Handler(mux),
	}

	return httpServer
}
