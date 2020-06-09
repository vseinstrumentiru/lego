package jaegerwrap

import (
	jaegerPropagation "contrib.go.opencensus.io/exporter/jaeger/propagation"
	"fmt"
	"go.opencensus.io/trace"
	"net/http"
	"strings"
)

const (
	httpHeader = `uber-trace-id`
)

var parent = &jaegerPropagation.HTTPFormat{}

type HTTPFormat struct{}

func (f *HTTPFormat) SpanContextFromRequest(req *http.Request) (sc trace.SpanContext, ok bool) {
	return parent.SpanContextFromRequest(req)
}

func (f *HTTPFormat) SpanContextToRequest(sc trace.SpanContext, req *http.Request) {
	header := fmt.Sprintf("%s:%s:%s:%d",
		strings.Replace(sc.TraceID.String(), "0000000000000000", "", 1), // Replacing 0 if string is 8bit
		sc.SpanID.String(),
		"0", // Parent span deprecated and will therefore be ignored.
		int64(sc.TraceOptions))
	req.Header.Set(httpHeader, header)
}
