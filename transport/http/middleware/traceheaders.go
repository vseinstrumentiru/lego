package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
)

type HeaderToTagsConfig map[string]string

func TraceHeadersMiddleware(cfg HeaderToTagsConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}

			for header, tag := range cfg {
				span.AddAttributes(trace.StringAttribute(tag, r.Header.Get(header)))
			}

			next.ServeHTTP(w, r)
		})
	}
}
