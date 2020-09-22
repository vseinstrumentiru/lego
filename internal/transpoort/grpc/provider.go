package http

import (
	"fmt"

	"emperror.dev/emperror"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/vseinstrumentiru/lego/internal/metrics/propagation"
	"github.com/vseinstrumentiru/lego/metrics/tracing"
	"github.com/vseinstrumentiru/lego/multilog"
	grpccfg "github.com/vseinstrumentiru/lego/transport/grpc"
	"github.com/vseinstrumentiru/lego/transport/http/miiddleware"
	"github.com/vseinstrumentiru/lego/version"
)

type args struct {
	dig.In
	Config      *grpccfg.Config
	TraceTags   miiddleware.TraceTagsMiddlewareConfig `optional:"true"`
	TraceConfig *tracing.Config                       `optional:"true"`

	Propagation *propagation.HTTPFormatCollection

	Version  version.Info
	Logger   multilog.Logger
	Pipeline *run.Group
	Upg      *tableflip.Upgrader
}

func Provide(in args) *grpc.Server {
	logger := in.Logger.WithFields(map[string]interface{}{"component": "grpc"})

	startOptions := trace.StartOptions{
		Sampler:  nil,
		SpanKind: trace.SpanKindServer,
	}

	if in.TraceConfig != nil {
		startOptions.Sampler = trace.Sampler(in.TraceConfig.Sampler)
	}

	server := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
		StartOptions:     startOptions,
		IsPublicEndpoint: in.Config.IsPublic,
	}))

	grpcLn, err := in.Upg.Listen("tcp", fmt.Sprintf(":%v", in.Config.Port))
	emperror.Panic(err)

	emperror.Panic(view.Register(
		ocgrpc.ServerReceivedBytesPerRPCView,
		ocgrpc.ServerSentBytesPerRPCView,
		ocgrpc.ServerLatencyView,
		ocgrpc.ServerCompletedRPCsView,
	))

	in.Pipeline.Add(appkitrun.LogServe(logger.WithFields(map[string]interface{}{"port": in.Config.Port}))(
		appkitrun.GRPCServe(server, grpcLn),
	))

	return server
}
