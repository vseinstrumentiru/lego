package httptools

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	"net/http"
)

type Handler interface {
	Handle(ctx context.Context, request interface{}) (response interface{}, err error)
	Decode(context.Context, *http.Request) (request interface{}, err error)
	Encode(context.Context, http.ResponseWriter, interface{}) error
	Middleware() []endpoint.Middleware
}

var errorEncoder = kitxhttp.NewJSONProblemErrorResponseEncoder(appkithttp.NewDefaultProblemConverter())

func JSONResponseEncoder(ctx context.Context, w http.ResponseWriter, resp interface{}, code int) error {
	return kitxhttp.JSONResponseEncoder(ctx, w, kitxhttp.WithStatusCode(resp, code))
}

func New(handler Handler, options []kithttp.ServerOption) http.Handler {
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
	dec kithttp.DecodeRequestFunc,
	enc kithttp.EncodeResponseFunc,
	middleware []endpoint.Middleware,
	options []kithttp.ServerOption,
) http.Handler {
	if len(middleware) > 0 {
		e = kitxendpoint.Combine(middleware...)(e)
	}

	return kithttp.NewServer(
		e,
		dec,
		kitxhttp.ErrorResponseEncoder(enc, errorEncoder),
		options...,
	)
}
