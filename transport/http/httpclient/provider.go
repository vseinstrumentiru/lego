package httpclient

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"go.uber.org/dig"

	propagationx "github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	httpTransport "github.com/vseinstrumentiru/lego/v2/transport/http"
)

type ClientArgs struct {
	dig.In
	Propagation *propagationx.HTTPFormatCollection `optional:"true"`
}

func ConstructorProvider(in ClientArgs) httpTransport.Constructor {
	var prop propagation.HTTPFormat

	if in.Propagation != nil {
		prop = in.Propagation
	}

	c := func(name string) *http.Client {
		return &http.Client{
			Transport: &ochttp.Transport{
				Base:        http.DefaultTransport,
				Propagation: prop,
				StartOptions: trace.StartOptions{
					Sampler: trace.AlwaysSample(),
				},
				FormatSpanName: func(req *http.Request) string {
					return name + ":" + req.URL.Path
				},
			},
		}
	}

	return c
}

func Provide(in ClientArgs) *http.Client {
	var prop propagation.HTTPFormat

	if in.Propagation != nil {
		prop = in.Propagation
	}

	c := &http.Client{
		Transport: &ochttp.Transport{
			Base:        http.DefaultTransport,
			Propagation: prop,
			StartOptions: trace.StartOptions{
				Sampler: trace.AlwaysSample(),
			},
			FormatSpanName: func(req *http.Request) string {
				return req.URL.Path
			},
		},
	}

	return c
}
