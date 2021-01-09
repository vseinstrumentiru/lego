package httpserver

import (
	"fmt"
	"net/http"

	"emperror.dev/emperror"
	"github.com/cloudflare/tableflip"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/sagikazarmark/ocmux"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	httpcfg "github.com/vseinstrumentiru/lego/v2/transport/http"
	"github.com/vseinstrumentiru/lego/v2/transport/http/middleware"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type Args struct {
	dig.In
	Config      *httpcfg.Config                      `optional:"true"`
	TraceTags   middleware.TraceTagsMiddlewareConfig `optional:"true"`
	TraceConfig *tracing.Config                      `optional:"true"`
	Newrelic    *newrelic.Application                `optional:"true"`

	Propagation *propagation.HTTPFormatCollection

	Version  version.Info
	Logger   multilog.Logger
	Pipeline *run.Group
	Upg      *tableflip.Upgrader
}

func Provide(in Args) (*http.Server, *mux.Router) {
	if in.Config == nil {
		in.Config = httpcfg.NewDefaultConfig()
	}

	logger := in.Logger.WithFields(map[string]interface{}{"component": "http"})

	router := mux.NewRouter()
	router.Use(middleware.RecoverHandlerMiddleware(logger))

	startOptions := trace.StartOptions{
		Sampler:  nil,
		SpanKind: trace.SpanKindServer,
	}

	if in.TraceConfig != nil {
		router.Use(ocmux.Middleware())
		router.Use(middleware.TraceVersionMiddleware(in.Version))
		if len(in.TraceTags) > 0 {
			router.Use(middleware.TraceTagsMiddleware(in.TraceTags))
		}

		startOptions.Sampler = trace.Sampler(in.TraceConfig.Sampler)
	}

	if in.Newrelic != nil {
		router.Use(nrgorilla.Middleware(in.Newrelic))
	}

	handler := &ochttp.Handler{
		Handler:          router,
		Propagation:      in.Propagation,
		StartOptions:     startOptions,
		IsPublicEndpoint: in.Config.IsPublic,
	}

	server := &http.Server{
		Handler:  handler,
		ErrorLog: multilog.NewErrorStandardLogger(logger),
	}

	httpLn, err := in.Upg.Listen("tcp", fmt.Sprintf(":%v", in.Config.Port))
	emperror.Panic(err)

	emperror.Panic(view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	))

	in.Pipeline.Add(appkitrun.LogServe(logger.WithFields(map[string]interface{}{"port": in.Config.Port}))(
		appkitrun.HTTPServe(server, httpLn, in.Config.ShutdownTimeout),
	))

	return server, router
}
