package grpc

import (
	"emperror.dev/emperror"
	appkitrun "github.com/sagikazarmark/appkit/run"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"io"
	"logur.dev/logur"
	"strconv"
)

func Run(p lego2.Process, config Config) (*grpc.Server, io.Closer) {
	const name = "grpc"

	logger := logur.WithField(p.Log(), "server", name)

	server := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
		StartOptions: trace.StartOptions{
			Sampler:  trace.AlwaysSample(),
			SpanKind: trace.SpanKindServer,
		},
		IsPublicEndpoint: config.IsPublic,
	}))

	addr := ":" + strconv.Itoa(config.Port)
	logger.Info("listening on address", map[string]interface{}{"address": addr})

	grpcLn, err := p.Listen("tcp", addr)
	emperror.Panic(err)

	p.Background(appkitrun.LogServe(logger)(appkitrun.GRPCServe(server, grpcLn)))

	return server, lego2.CloseFn(func() error {
		server.Stop()
		return nil
	})
}
