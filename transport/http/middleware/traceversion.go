package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"

	"github.com/vseinstrumentiru/lego/v2/version"
)

func TraceVersionMiddleware(buildInfo version.Info) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}

			span.AddAttributes(
				trace.StringAttribute("build.version", buildInfo.Version),
				trace.StringAttribute("build.commit", buildInfo.CommitHash),
				trace.StringAttribute("build.date", buildInfo.BuildDate),
			)
			next.ServeHTTP(w, r)
		})
	}
}
