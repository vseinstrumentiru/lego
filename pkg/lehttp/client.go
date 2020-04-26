package lehttp

import (
	"crypto/tls"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"github.com/vseinstrumentiru/lego/internal/monitor/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"net/url"
)

type ClientOption func(*http.Client)

func SetSpanNameFormatter(formatter func(req *http.Request) string) ClientOption {
	return func(client *http.Client) {
		if octr, ok := client.Transport.(*ochttp.Transport); ok {
			octr.FormatSpanName = formatter
		}
	}
}

func SetTraceSampler(sampler trace.Sampler) ClientOption {
	return func(client *http.Client) {
		if octr, ok := client.Transport.(*ochttp.Transport); ok {
			octr.StartOptions.Sampler = sampler
		}
	}
}

func InsecureClient() ClientOption {
	return func(client *http.Client) {
		var httpTr *http.Transport
		if octr, ok := client.Transport.(*ochttp.Transport); ok {
			if octr.Base == nil {
				octr.Base = http.DefaultTransport
			}

			httpTr = octr.Base.(*http.Transport)
		} else if httpTr, ok = client.Transport.(*http.Transport); !ok {
			return
		}

		if httpTr == nil {
			return
		}

		if httpTr.TLSClientConfig == nil {
			httpTr.TLSClientConfig = &tls.Config{}
		}

		httpTr.TLSClientConfig.InsecureSkipVerify = true
	}
}

func NewClient(name string, opts ...ClientOption) *http.Client {
	c := &http.Client{
		Transport: &ochttp.Transport{
			Base:        http.DefaultTransport,
			Propagation: propagation.DefaultHTTPFormat,
			StartOptions: trace.StartOptions{
				Sampler: trace.AlwaysSample(),
			},
			FormatSpanName: func(req *http.Request) string {
				return name + ":" + req.URL.Path
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func NewRestyClient(name string, host string, opts ...ClientOption) *resty.Client {
	return resty.NewWithClient(NewClient(name, opts...)).SetHostURL(host)
}

func NewGraphQLClient(name string, endpoint string, opts ...ClientOption) *graphql.Client {
	return graphql.NewClient(endpoint, NewClient(name))
}

func NewJsonRPCClient(addr, method string, client *http.Client, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	u, err := url.Parse(addr)
	emperror.Panic(errors.Wrap(err, "can't create JsonRPC client"))

	option = append([]jsonrpc.ClientOption{jsonrpc.SetClient(client)}, option...)

	return jsonrpc.NewClient(u, method, option...)
}
