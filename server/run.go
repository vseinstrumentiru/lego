package server

import (
	lego "github.com/vseinstrumentiru/lego/v2/app"
	"github.com/vseinstrumentiru/lego/v2/module"
	"github.com/vseinstrumentiru/lego/v2/transport/mysql"
)

func Run(app interface{}, opts ...lego.Option) {
	opts = append([]lego.Option{
		lego.ServerMode(),
		lego.Provide(
			module.All,
			mysql.ProvideConnector,
		),
	}, opts...)
	lego.NewRuntime(opts...).Run(app)
}
