package lehttp

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
	"net/url"
)

func NewClient() *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{},
	}
}

func NewRestyClient(host string) *resty.Client {
	return resty.NewWithClient(NewClient()).SetHostURL(host)
}

func NewGraphQLClient(endpoint string) *graphql.Client {
	return graphql.NewClient(endpoint, NewClient())
}

func NewJsonRPCClient(addr, method string, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	u, err := url.Parse(addr)
	emperror.Panic(errors.Wrap(err, "can't create JsonRPC client"))

	option = append([]jsonrpc.ClientOption{jsonrpc.SetClient(NewClient())}, option...)

	return jsonrpc.NewClient(u, method, option...)
}
