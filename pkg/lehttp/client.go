package lehttp

import (
	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
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
