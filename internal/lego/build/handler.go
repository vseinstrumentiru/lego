package build

import (
	"encoding/json"
	"go.opencensus.io/trace"
	"net/http"

	"emperror.dev/errors"
)

// Handler returns an HTTP handler for version information.
func Handler(buildInfo Info) http.Handler {
	var body []byte

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body == nil {
			var err error

			body, err = json.Marshal(buildInfo)
			if err != nil {
				panic(errors.Wrap(err, "failed to render version information"))
			}
		}

		_, _ = w.Write(body)
	})
}

func TraceMiddleware(buildInfo Info) func(http.Handler) http.Handler {
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
