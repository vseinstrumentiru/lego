package metrics

import (
	"fmt"
	"net/http"

	"emperror.dev/emperror"
	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type args struct {
	dig.In
	Config *metrics.Config `optional:"true"`
	Trace  *tracing.Config `optional:"true"`

	Version  version.Info
	Logger   multilog.Logger
	Pipeline *run.Group
	Upg      *tableflip.Upgrader
}

func Provide(in args) (*http.ServeMux, health.Health) {
	server := http.DefaultServeMux
	server.Handle("/version", version.Handler(in.Version))

	logger := in.Logger.WithFields(map[string]interface{}{"component": "metrics"})

	if in.Config == nil {
		in.Config = metrics.NewDefaultConfig()
	}

	if in.Config.Debug {
		zpages.Handle(server, "/debug")
	}

	healthz := health.New()
	healthz.WithCheckListener(NewLogCheckListener(logger))

	{
		handler := healthhttp.HandleHealthJSON(healthz)
		server.Handle("/healthz", handler)

		// Kubernetes style health checks
		server.HandleFunc("/healthz/live", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
		server.Handle("/healthz/ready", handler)
	}

	emperror.Panic(view.Register(
		// Health checks
		health.ViewCheckCountByNameAndStatus,
		health.ViewCheckStatusByName,
		health.ViewCheckExecutionTime,
		// HTTP Client
		ochttp.ClientCompletedCount,
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		// GRPC Client
		ocgrpc.ClientSentBytesPerRPCView,
		ocgrpc.ClientReceivedBytesPerRPCView,
		ocgrpc.ClientRoundtripLatencyView,
		ocgrpc.ClientRoundtripLatencyView,
	))

	httpLn, err := in.Upg.Listen("tcp", fmt.Sprintf(":%v", in.Config.Port))
	emperror.Panic(err)

	srv := &http.Server{
		Handler:  server,
		ErrorLog: multilog.NewErrorStandardLogger(logger),
	}

	in.Pipeline.Add(appkitrun.LogServe(logger.WithFields(map[string]interface{}{"port": in.Config.Port}))(
		appkitrun.HTTPServe(srv, httpLn, 0),
	))

	if in.Trace != nil {
		trace.ApplyConfig(
			trace.Config{
				DefaultSampler:             trace.Sampler(in.Trace.Sampler),
				MaxAnnotationEventsPerSpan: in.Trace.MaxAnnotationEventsPerSpan,
				MaxMessageEventsPerSpan:    in.Trace.MaxMessageEventsPerSpan,
				MaxAttributesPerSpan:       in.Trace.MaxAttributesPerSpan,
				MaxLinksPerSpan:            in.Trace.MaxLinksPerSpan,
			},
		)
	}

	return server, healthz
}
