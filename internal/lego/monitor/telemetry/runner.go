package telemetry

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	health "github.com/AppsFlyer/go-sundheit"
	appkitrun "github.com/sagikazarmark/appkit/run"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/log"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"io"
	"logur.dev/logur"
	"net/http"
)

func Run(p lego2.Process, router *http.ServeMux, config Config) io.Closer {
	const name = "telemetry"
	logger := logur.WithField(p.Log(), "server", name)

	logger.Info("listening on address", map[string]interface{}{"address": config.Addr})

	ln, err := p.Listen("tcp", config.Addr)
	emperror.Panic(err)

	server := &http.Server{
		Handler:  router,
		ErrorLog: log.NewErrorStandardLogger(logger),
	}

	p.Background(appkitrun.LogServe(logger)(appkitrun.HTTPServe(server, ln, p.ShutdownTimeout())))

	registerStats(config.Stats)

	return server
}

func registerStats(stats []*view.View) {
	err := view.Register(
		// Health checks
		health.ViewCheckCountByNameAndStatus,
		health.ViewCheckStatusByName,
		health.ViewCheckExecutionTime,

		// HTTP Client
		ochttp.ClientCompletedCount,
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,

		// HTTP
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,

		// GRPC Client
		ocgrpc.ClientSentBytesPerRPCView,
		ocgrpc.ClientReceivedBytesPerRPCView,
		ocgrpc.ClientRoundtripLatencyView,
		ocgrpc.ClientRoundtripLatencyView,

		// GRPC
		ocgrpc.ServerReceivedBytesPerRPCView,
		ocgrpc.ServerSentBytesPerRPCView,
		ocgrpc.ServerLatencyView,
		ocgrpc.ServerCompletedRPCsView,
	)

	emperror.Panic(errors.Wrap(err, "failed to register stat views"))

	if len(stats) > 0 {
		err = view.Register(stats...)

		emperror.Panic(errors.Wrap(err, "failed to register app stat views"))
	}
}
