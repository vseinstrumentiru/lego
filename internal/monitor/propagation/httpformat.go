package propagation

import (
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"net/http"
)

var DefaultHTTPFormat propagation.HTTPFormat

type HTTPFormatCollection []propagation.HTTPFormat

func (c HTTPFormatCollection) SpanContextFromRequest(req *http.Request) (sc trace.SpanContext, ok bool) {
	for _, f := range c {
		if sc, ok := f.SpanContextFromRequest(req); ok {
			return sc, ok
		}
	}

	return trace.SpanContext{}, false
}

func (c HTTPFormatCollection) SpanContextToRequest(sc trace.SpanContext, req *http.Request) {
	for _, f := range c {
		f.SpanContextToRequest(sc, req)
	}
}
