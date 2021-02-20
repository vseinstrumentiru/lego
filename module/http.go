package module

import (
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpclient"
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpserver"
)

func HTTPServer() (interface{}, []interface{}) {
	return httpserver.ProvideServer, nil
}

func HTTPClient() (interface{}, []interface{}) {
	return httpclient.Provide, nil
}

func HTTPClientConstructor() (interface{}, []interface{}) {
	return httpclient.ConstructorProvider, nil
}

func MuxRouter() (interface{}, []interface{}) {
	return httpserver.ProvideMuxRouter, nil
}
