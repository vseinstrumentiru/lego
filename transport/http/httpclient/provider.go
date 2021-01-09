package httpclient

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	httpTransport "github.com/vseinstrumentiru/lego/v2/transport/http"
)

func ConstructorProvider(httpProp *propagation.HTTPFormatCollection) httpTransport.Constructor {
	return func(name string) *http.Client {
		return &http.Client{
			Transport: &ochttp.Transport{
				Base:        http.DefaultTransport,
				Propagation: httpProp,
				StartOptions: trace.StartOptions{
					Sampler: trace.AlwaysSample(),
				},
				FormatSpanName: func(req *http.Request) string {
					return name + ":" + req.URL.Path
				},
			},
		}
	}
}

func Provide(httpProp *propagation.HTTPFormatCollection) *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{
			Base:        http.DefaultTransport,
			Propagation: httpProp,
			StartOptions: trace.StartOptions{
				Sampler: trace.AlwaysSample(),
			},
			FormatSpanName: func(req *http.Request) string {
				return req.URL.Path
			},
		},
	}
}
