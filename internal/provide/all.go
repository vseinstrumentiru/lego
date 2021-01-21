package provide

import "github.com/vseinstrumentiru/lego/v2/config"

func All(runtime config.Runtime) []interface{} {
	if runtime.Is(config.OptWithoutProviders) {
		return []interface{}{}
	}

	var res []interface{}

	res = append(res, Minimal()...)

	if runtime.Is(config.ServerMode) {
		res = append(res, Monitoring()...)
		res = append(res, Pipeline()...)
	}

	res = append(res, Http()...)
	res = append(res, Grpc()...)
	res = append(res, Sql()...)
	res = append(res, Events()...)

	return res
}
