package http

import (
	"emperror.dev/emperror"
	"github.com/gorilla/mux"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/sagikazarmark/ocmux"
	"github.com/vseinstrumentiru/lego/internal/monitor/log"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"io"
	"logur.dev/logur"
	"net/http"
)

func Run(p lego.Process, config Config) (*mux.Router, io.Closer) {
	const name = "http"

	logger := logur.WithField(p.Log(), "server", name)

	router := mux.NewRouter()
	router.Use(ocmux.Middleware())

	server := &http.Server{
		Handler: &ochttp.Handler{
			Handler: router,
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
			IsPublicEndpoint: true,
		},
		ErrorLog: log.NewErrorStandardLogger(logger),
	}

	logger.Info("listening on address", map[string]interface{}{"address": config.Addr})

	httpLn, err := p.Listen("tcp", config.Addr)
	emperror.Panic(err)

	p.Run(appkitrun.LogServe(logger)(appkitrun.HTTPServe(server, httpLn, p.ShutdownTimeout())))

	return router, server
}
