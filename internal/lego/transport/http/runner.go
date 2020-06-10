package http

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/gorilla/mux"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/sagikazarmark/ocmux"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/build"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/log"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"io"
	"logur.dev/logur"
	"net/http"
	"strconv"
)

func Run(p lego2.Process, config Config) (*mux.Router, io.Closer) {
	const name = "http"

	logger := logur.WithField(p.Log(), "server", name)

	router := mux.NewRouter()
	router.Use(recoverHandlerMiddleware(p.Handle), ocmux.Middleware())

	traceCfg := additionalTagsConfig{DataCenter: p.DataCenterName()}

	server := &http.Server{
		Handler: &ochttp.Handler{
			Handler:     additionalTagsMiddleware(traceCfg)(build.TraceMiddleware(p.Build())(router)),
			Propagation: propagation.DefaultHTTPFormat,
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
			IsPublicEndpoint: config.IsPublic,
		},
		ErrorLog: log.NewErrorStandardLogger(logger),
	}

	addr := ":" + strconv.Itoa(config.Port)
	logger.Info("listening on address", map[string]interface{}{"address": addr})

	httpLn, err := p.Listen("tcp", addr)
	emperror.Panic(err)

	p.Background(appkitrun.LogServe(logger)(appkitrun.HTTPServe(server, httpLn, p.ShutdownTimeout())))

	return router, server
}

type additionalTagsConfig struct {
	DataCenter string
}

func additionalTagsMiddleware(cfg additionalTagsConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}

			span.AddAttributes(
				trace.StringAttribute("server.dc", cfg.DataCenter),
			)
			next.ServeHTTP(w, r)
		})
	}
}

func recoverHandlerMiddleware(errorHandler func(err error)) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					errorHandler(errors.Wrap(err.(error), "http server panic"))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}