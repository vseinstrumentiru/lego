package http

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type Constructor func(name string) *http.Client

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

			//nolint:forcetypeassert
			httpTr = octr.Base.(*http.Transport)
		} else if httpTr, ok = client.Transport.(*http.Transport); !ok {
			return
		}

		if httpTr == nil {
			return
		}

		if httpTr.TLSClientConfig == nil {
			httpTr.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}

		httpTr.TLSClientConfig.InsecureSkipVerify = true
	}
}

func NewRestyClient(client *http.Client, host string) *resty.Client {
	return resty.NewWithClient(client).SetHostURL(host)
}

func NewGraphQLClient(client *http.Client, endpoint string) *graphql.Client {
	return graphql.NewClient(endpoint, client)
}

func NewJSONRPCClient(addr, method string, client *http.Client, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	u, err := url.Parse(addr)
	emperror.Panic(errors.Wrap(err, "can't create JsonRPC client"))

	option = append([]jsonrpc.ClientOption{jsonrpc.SetClient(client)}, option...)

	return jsonrpc.NewClient(u, method, option...)
}
