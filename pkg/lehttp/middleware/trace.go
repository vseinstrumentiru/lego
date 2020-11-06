package middleware

import (
	"github.com/gorilla/mux"

	"github.com/vseinstrumentiru/lego/transport/http/middleware"
)

// Deprecated: use middleware.LogRequestWithMaxLenMiddleware, middleware.LogResponseWithMaxLenMiddleware
// middleware.TraceHeadersMiddleware
type TraceRequestOptions struct {
	LogRequest        bool
	RequestBodyLimit  int
	LogResponse       bool
	ResponseBodyLimit int
	HeadersToTags     map[string]string
}

// Deprecated: use middleware.LogRequestWithMaxLenMiddleware, middleware.LogResponseWithMaxLenMiddleware
// middleware.TraceHeadersMiddleware
func TraceRequestResponse(opt TraceRequestOptions) mux.MiddlewareFunc {
	var mv []mux.MiddlewareFunc

	if opt.LogRequest {
		mv = append(mv, middleware.LogRequestWithMaxLenMiddleware(opt.RequestBodyLimit))
	}

	if opt.LogResponse {
		mv = append(mv, middleware.LogResponseWithMaxLenMiddleware(opt.ResponseBodyLimit))
	}

	if len(opt.HeadersToTags) > 0 {
		mv = append(mv, middleware.TraceHeadersMiddleware(opt.HeadersToTags))
	}

	return middleware.Combine(mv...)
}
