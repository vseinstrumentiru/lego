package middleware

import (
	"github.com/vseinstrumentiru/lego/tools/lehttp/middleware"
	"net/http"
)

// deprecated
func TraceRequestResponse(opt middleware.TraceRequestOptions) func(http.Handler) http.Handler {
	return middleware.TraceRequestResponse(opt)
}

// deprecated
func TraceTagFromHeaders(headersTagNames map[string]string) func(http.Handler) http.Handler {
	return middleware.TraceTagFromHeaders(headersTagNames)
}
