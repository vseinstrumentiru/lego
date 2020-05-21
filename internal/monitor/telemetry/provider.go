package telemetry

import (
	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/vseinstrumentiru/lego/pkg/build"
	"go.opencensus.io/zpages"
	"logur.dev/logur"
	"net/http"
)

func Provide(logger logur.LoggerFacade, info build.Info) (*http.ServeMux, health.Health) {
	telemetryRouter := http.DefaultServeMux
	telemetryRouter.Handle("/buildinfo", build.Handler(info))
	zpages.Handle(telemetryRouter, "/debug")

	healthz := health.New()
	healthz.WithCheckListener(NewLogger(logur.WithField(logger, "component", "healthcheck")))

	{
		handler := healthhttp.HandleHealthJSON(healthz)
		telemetryRouter.Handle("/healthz", handler)

		// Kubernetes style health checks
		telemetryRouter.HandleFunc("/healthz/live", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
		telemetryRouter.Handle("/healthz/ready", handler)

		return telemetryRouter, healthz
	}
}
