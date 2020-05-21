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
	LogRequest        bool
	RequestBodyLimit  int
	LogResponse       bool
	ResponseBodyLimit int
	HeadersToTags     map[string]string
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

				if opt.RequestBodyLimit > 0 && len(bodyBytes) > opt.RequestBodyLimit {
					bodyBytes = bodyBytes[:opt.RequestBodyLimit-1]
				}

				reqLog := string(bodyBytes)

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

			var bodyBytes []byte
			bodyBytes = tw.buf
			if opt.ResponseBodyLimit > 0 && len(bodyBytes) > opt.ResponseBodyLimit {
				bodyBytes = bodyBytes[:opt.ResponseBodyLimit-1]
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
