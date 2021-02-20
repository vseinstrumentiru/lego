package grpc

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

	"github.com/vseinstrumentiru/lego/v2/common/netx"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/transport/http/middleware"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type Args struct {
	dig.In
	Config      *Config                              `optional:"true"`
	TraceTags   middleware.TraceTagsMiddlewareConfig `optional:"true"`
	TraceConfig *tracing.Config                      `optional:"true"`

	Propagation *propagation.HTTPFormatCollection

	Version  *version.Info
	Logger   multilog.Logger
	Pipeline *run.Group
	Upg      *tableflip.Upgrader `optional:"true"`
}

func Provide(in Args) *grpc.Server {
	if in.Config == nil {
		in.Config = NewDefaultConfig()
	}

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

	grpcLn, err := netx.Listen("tcp", fmt.Sprintf(":%v", in.Config.Port), in.Upg)
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
