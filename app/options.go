package app

import (
	config "github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/di"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
)

type Option func(r *runtime)

func resolveModulePack(pack []di.Module) (providers []interface{}, configurations []interface{}) {
	for _, m := range pack {
		p, c := m()
		if p != nil {
			providers = append(providers, p)
		}

		if len(c) > 0 {
			configurations = append(configurations, c...)
		}
	}

	return
}

func Provide(providers ...interface{}) Option {
	return func(r *runtime) {
		for _, provider := range providers {
			if pack, ok := provider.(func() []di.Module); ok {
				p, c := resolveModulePack(pack())
				for _, i := range p {
					if i != nil {
						r.container.Register(i)
					}
				}
				r.configurations = append(r.configurations, c...)
			} else if m, ok := provider.(func() (interface{}, []interface{})); ok {
				p, c := m()
				r.configurations = append(r.configurations, c...)
				if p != nil {
					r.container.Register(p)
				}
			} else {
				r.container.Register(provider)
			}
		}
	}
}

func LocalDebug() Option {
	return func(r *runtime) {
		r.opts.Set(config.OptLocalDebug, true)
	}
}

// Deprecated: use app.NewRuntime().Run()
func NoDefaultProviders() Option {
	return func(r *runtime) {}
}

func CommandMode() Option {
	return func(r *runtime) {
		r.exec = command
	}
}

func ServerMode() Option {
	return func(r *runtime) {
		r.exec = serve
	}
}

// Deprecated: use app.CommandMode()
func NoWait() Option {
	return CommandMode()
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
