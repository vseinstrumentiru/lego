package app

import (
	config "github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
)

type Option func(r *runtime)

func Provide(providers ...interface{}) Option {
	return func(r *runtime) {
		for _, provider := range providers {
			r.container.Register(provider)
		}
	}
}

func LocalDebug() Option {
	return func(r *runtime) {
		r.opts.Set(config.OptLocalDebug, true)
	}
}

func NoDefaultProviders() Option {
	return func(r *runtime) {
		r.opts.Set(config.OptWithoutProviders, true)
	}
}

func CommandMode() Option {
	return func(r *runtime) {
		r.exec = command
		r.opts.Set(config.ServerMode, false)
	}
}

func ServerMode() Option {
	return func(r *runtime) {
		r.opts.Set(config.ServerMode, true)
		r.exec = serve
	}
}

func NoWait() Option {
	return func(r *runtime) {
		r.opts.Set(config.ServerMode, false)
	}
}

func EnvPath(path string) Option {
	return func(r *runtime) {
		r.opts.Set(config.OptEnvPath, path)
	}
}

func WithConfig(cfg interface{}) Option {
	return func(r *runtime) {
		r.container.Register(func() env.Config {
			return cfg
		})
	}
}
