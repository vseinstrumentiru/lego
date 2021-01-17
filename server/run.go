package server

import (
	lego "github.com/vseinstrumentiru/lego/v2/app"
)

func Run(app interface{}, opts ...lego.Option) {
	opts = append([]lego.Option{lego.ServerMode()}, opts...)
	lego.NewRuntime(opts...).Run(app)
}
