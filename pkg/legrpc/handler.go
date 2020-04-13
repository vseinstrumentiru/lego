package legrpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
)

type Handler interface {
	Handle(ctx context.Context, request interface{}) (response interface{}, err error)
	Decode(context.Context, interface{}) (request interface{}, err error)
	Encode(context.Context, interface{}) (response interface{}, err error)
	Middleware() []endpoint.Middleware
}

var errorEncoder = kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

func New(handler Handler, options []kitgrpc.ServerOption) kitgrpc.Handler {
	return NewHandler(
		handler.Handle,
		handler.Decode,
		handler.Encode,
		handler.Middleware(),
		options,
	)
}

func NewHandler(
	e endpoint.Endpoint,
	dec kitgrpc.DecodeRequestFunc,
	enc kitgrpc.EncodeResponseFunc,
	middleware []endpoint.Middleware,
	options []kitgrpc.ServerOption,
) kitgrpc.Handler {
	if len(middleware) > 0 {
		e = kitxendpoint.Combine(middleware...)(e)
	}

	return kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
		e,
		dec,
		kitxgrpc.ErrorResponseEncoder(enc, errorEncoder),
		options...,
	), errorEncoder)
}
