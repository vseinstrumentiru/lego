package server

import (
	"github.com/vseinstrumentiru/lego/v2/app"
)

type Option = app.Option

func Provide(providers ...interface{}) Option {
	return app.Provide(providers...)
}

func LocalDebug() Option {
	return app.LocalDebug()
}

func NoDefaultProviders() Option {
	return app.NoDefaultProviders()
}

func CommandMode() Option {
	return app.CommandMode()
}

func NoWait() Option {
	return app.NoWait()
}

func EnvPath(path string) Option {
	return app.EnvPath(path)
}

func WithConfig(cfg interface{}) Option {
	return app.WithConfig(cfg)
}
