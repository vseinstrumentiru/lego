package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
)

type traceWriter struct {
	http.ResponseWriter
	buf []byte
}

func (w *traceWriter) Write(b []byte) (int, error) {
	i, err := w.ResponseWriter.Write(b)

	if err == nil {
		w.buf = []byte(string(w.buf) + string(b))
	}

	return i, err
}

func LogResponseMiddleware(length int) mux.MiddlewareFunc {
	return LogResponseWithMaxLenMiddleware(1024)
}

func LogResponseWithMaxLenMiddleware(length int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)

				return
			}

			tw := &traceWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(tw, r)

			var bodyBytes []byte
			bodyBytes = tw.buf
			if length > 0 && len(bodyBytes) > length {
				bodyBytes = bodyBytes[:length-1]
			}

			span.Annotate(
				[]trace.Attribute{
					trace.StringAttribute("response", string(bodyBytes)),
				},
				"response log",
			)
		})
	}
}
