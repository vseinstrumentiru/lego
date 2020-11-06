package lehttp

import (
	baseHttp "net/http"

	"emperror.dev/emperror"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"

	"github.com/vseinstrumentiru/lego/v2/internal/deprecated"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/http/httpclient"
	httpTransport "github.com/vseinstrumentiru/lego/v2/transport/http"
)

type ClientOption = httpTransport.ClientOption

// Deprecated: use DI with httpclient.Constructor
func NewClient(name string, opts ...ClientOption) *baseHttp.Client {
	var client *baseHttp.Client
	err := deprecated.Container.Execute(func(c httpclient.Constructor) {
		client = c(name)
	})

	emperror.Panic(err)

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Deprecated: use DI with httpclient.Constructor and http.NewRestyClient
func NewRestyClient(name string, host string, opts ...ClientOption) *resty.Client {
	c := NewClient(name, opts...)

	return httpTransport.NewRestyClient(c, host)
}

// Deprecated: use DI with httpclient.Constructor and http.NewGraphQLClient
func NewGraphQLClient(name string, endpoint string, opts ...ClientOption) *graphql.Client {
	c := NewClient(name, opts...)

	return httpTransport.NewGraphQLClient(c, endpoint)
}

// Deprecated: use http.NewJSONRPCClient
func NewJsonRPCClient(addr, method string, client *baseHttp.Client, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	return httpTransport.NewJSONRPCClient(addr, method, client, option...)
}
