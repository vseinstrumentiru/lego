package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
)

type TraceTagsMiddlewareConfig map[string]string

func TraceTagsMiddleware(cfg TraceTagsMiddlewareConfig) mux.MiddlewareFunc {
	tags := make([]trace.Attribute, 0, len(cfg))

	for key, val := range cfg {
		tags = append(tags, trace.StringAttribute(key, val))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil || len(tags) == 0 {
				next.ServeHTTP(w, r)

				return
			}

			span.AddAttributes(tags...)
			next.ServeHTTP(w, r)
		})
	}
}
