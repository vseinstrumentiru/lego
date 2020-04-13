package grpc

import (
	"emperror.dev/emperror"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"io"
	"logur.dev/logur"
)

func Run(p lego.Process, config Config) (*grpc.Server, io.Closer) {
	const name = "grpc"

	logger := logur.WithField(p.Log(), "server", name)

	server := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
		StartOptions: trace.StartOptions{
			Sampler:  trace.AlwaysSample(),
			SpanKind: trace.SpanKindServer,
		},
		IsPublicEndpoint: true,
	}))

	logger.Info("listening on address", map[string]interface{}{"address": config.Addr})

	grpcLn, err := p.Listen("tcp", config.Addr)
	emperror.Panic(err)

	p.Run(appkitrun.LogServe(logger)(appkitrun.GRPCServe(server, grpcLn)))

	return server, lego.CloseFn(func() error {
		server.Stop()
		return nil
	})
}
