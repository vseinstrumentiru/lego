package propagation

import (
	"net/http"

	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

func ProvideHTTP() *HTTPFormatCollection {
	c := new(HTTPFormatCollection)
	return c
}

type HTTPFormatCollection struct {
	formatters []propagation.HTTPFormat
}

func (c *HTTPFormatCollection) Add(formatter propagation.HTTPFormat) *HTTPFormatCollection {
	c.formatters = append(c.formatters, formatter)

	return c
}

func (c *HTTPFormatCollection) SpanContextFromRequest(req *http.Request) (sc trace.SpanContext, ok bool) {
	for _, f := range c.formatters {
		if sc, ok := f.SpanContextFromRequest(req); ok {
			return sc, ok
		}
	}

	return trace.SpanContext{}, false
}

func (c *HTTPFormatCollection) SpanContextToRequest(sc trace.SpanContext, req *http.Request) {
	for _, f := range c.formatters {
		f.SpanContextToRequest(sc, req)
	}
}
