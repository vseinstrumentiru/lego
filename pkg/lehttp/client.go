package lehttp

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"net/url"
)

func NewClient(name string) *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{
			StartOptions: trace.StartOptions{
				Sampler: trace.AlwaysSample(),
			},
			FormatSpanName: func(req *http.Request) string {
				return name + ":" + req.URL.Path
			},
		},
	}
}

func NewRestyClient(name string, host string) *resty.Client {
	return resty.NewWithClient(NewClient(name)).SetHostURL(host)
}

func NewGraphQLClient(name string, endpoint string) *graphql.Client {
	return graphql.NewClient(endpoint, NewClient(name))
}

func NewJsonRPCClient(name string, addr, method string, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	u, err := url.Parse(addr)
	emperror.Panic(errors.Wrap(err, "can't create JsonRPC client"))

	option = append([]jsonrpc.ClientOption{jsonrpc.SetClient(NewClient(name))}, option...)

	return jsonrpc.NewClient(u, method, option...)
}
