package http

import (
	"net/http"
	"net/url"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
)

func NewRestyClient(client *http.Client, host string) *resty.Client {
	return resty.NewWithClient(client).SetHostURL(host)
}

func NewGraphQLClient(client *http.Client, endpoint string) *graphql.Client {
	return graphql.NewClient(endpoint, client)
}

func NewJsonRPCClient(addr, method string, client *http.Client, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	u, err := url.Parse(addr)
	emperror.Panic(errors.Wrap(err, "can't create JsonRPC client"))

	option = append([]jsonrpc.ClientOption{jsonrpc.SetClient(client)}, option...)

	return jsonrpc.NewClient(u, method, option...)
}
