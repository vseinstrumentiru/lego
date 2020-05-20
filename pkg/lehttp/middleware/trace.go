package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
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

type TraceRequestOptions struct {
	LogRequest    bool
	LogResponse   bool
	HeadersToTags map[string]string
}

func TraceRequestResponse(opt TraceRequestOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}

			if opt.HeadersToTags != nil {
				for header, tag := range opt.HeadersToTags {
					span.AddAttributes(trace.StringAttribute(tag, r.Header.Get(header)))
				}
			}

			if opt.LogRequest {
				var reqLog string
				if r.Method == http.MethodGet {
					req := struct {
						Query string            `json:"query"`
						Vars  map[string]string `json:"vars"`
					}{
						Query: r.URL.RawQuery,
						Vars:  mux.Vars(r),
					}
					bodyBytes, _ := json.Marshal(req)
					reqLog = string(bodyBytes)
				} else {
					bodyBytes, _ := ioutil.ReadAll(r.Body)
					_ = r.Body.Close() //  must close
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
					reqLog = string(bodyBytes)
				}

				span.Annotate(
					[]trace.Attribute{
						trace.StringAttribute("request", reqLog),
					},
					"request log",
				)
			}

			if !opt.LogResponse {
				next.ServeHTTP(w, r)
				return
			}

			tw := &traceWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(tw, r)

			span.Annotate(
				[]trace.Attribute{
					trace.StringAttribute("response", string(tw.buf)),
				},
				"response log",
			)
		})
	}
}

func TraceTagFromHeaders(headersTagNames map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.FromContext(r.Context())

			if span == nil {
				next.ServeHTTP(w, r)
				return
			}
		})
	}
}
