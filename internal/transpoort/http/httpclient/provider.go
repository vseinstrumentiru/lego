package httpclient

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func ProvideClient() *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{
			Base: http.DefaultTransport,
			// Propagation: propagation.DefaultHTTPFormat,
			StartOptions: trace.StartOptions{
				Sampler: trace.AlwaysSample(),
			},
		},
	}
}
