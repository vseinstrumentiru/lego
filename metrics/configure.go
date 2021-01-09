package metrics

import (
	health "github.com/AppsFlyer/go-sundheit"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
)

func ConfigureStats() error {
	return view.Register(
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
	)
}

type TraceConfigArgs struct {
	dig.In
	Trace *tracing.Config `optional:"true"`
}

func ConfigureTrace(in TraceConfigArgs) error {
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

	return nil
}
