package lehttp

import (
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"github.com/vseinstrumentiru/lego/tools/lehttp"
	"go.opencensus.io/trace"
	"net/http"
)

// deprecated
type ClientOption = httptools.ClientOption

// deprecated
func SetSpanNameFormatter(formatter func(req *http.Request) string) ClientOption {
	return httptools.SetSpanNameFormatter(formatter)
}

// deprecated
func SetTraceSampler(sampler trace.Sampler) ClientOption {
	return httptools.SetTraceSampler(sampler)
}

// deprecated
func InsecureClient() ClientOption {
	return httptools.InsecureClient()
}

// deprecated
func NewClient(name string, opts ...ClientOption) *http.Client {
	return httptools.NewClient(name, opts...)
}

// deprecated
func NewRestyClient(name string, host string, opts ...ClientOption) *resty.Client {
	return httptools.NewRestyClient(name, host, opts...)
}

// deprecated
func NewGraphQLClient(name string, endpoint string, opts ...ClientOption) *graphql.Client {
	return httptools.NewGraphQLClient(name, endpoint, opts...)
}

// deprecated
func NewJsonRPCClient(addr, method string, client *http.Client, option ...jsonrpc.ClientOption) *jsonrpc.Client {
	return httptools.NewJsonRPCClient(addr, method, client, option...)
}
