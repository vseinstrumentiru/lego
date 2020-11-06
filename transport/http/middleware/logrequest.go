package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
)

func LogRequestMiddleware() mux.MiddlewareFunc {
	return LogRequestWithMaxLenMiddleware(1024)
}

func LogRequestWithMaxLenMiddleware(length int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}

			var bodyBytes []byte
			if r.Method == http.MethodGet {
				req := struct {
					Query string            `json:"query"`
					Vars  map[string]string `json:"vars"`
				}{
					Query: r.URL.RawQuery,
					Vars:  mux.Vars(r),
				}
				bodyBytes, _ = json.Marshal(req)

			} else {
				bodyBytes, _ = ioutil.ReadAll(r.Body)
				_ = r.Body.Close() //  must close
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			if length > 0 && len(bodyBytes) > length {
				bodyBytes = bodyBytes[:length-1]
			}

			reqLog := string(bodyBytes)

			span.Annotate(
				[]trace.Attribute{
					trace.StringAttribute("request", reqLog),
				},
				"request log",
			)

			next.ServeHTTP(w, r)
		})
	}
}
