package server

const (
	optNoWait  = "opt_no_wait"
	optEnvPath = "opt_env_path"

	optLocalDebug = "opt_local_debug"

	optWithoutProviders = "opt_without_providers"

	optUseJaeger     = "opt_use_jaeger"
	optUsePrometheus = "opt_use_prometheus"
	optUseOpencensus = "opt_use_opencensus"
	optUseNewRelic   = "opt_use_newrelic"
	optUseTrace      = "opt_use_trace"
	optUseStats      = "opt_use_stats"
	optUseMonitoring = "opt_use_monitoring"
)

func LocalDebug() Option {
	return func(r *Runtime) {
		r.opts.Set(optLocalDebug, true)
	}
}

func NoDefaultProviders() Option {
	return func(r *Runtime) {
		r.opts.Set(optWithoutProviders, true)
	}
}

func NoWait() Option {
	return func(r *Runtime) {
		r.opts.Set(optNoWait, true)
	}
}

// Deprecated: use NoWait
func NoWaitOption() Option {
	return NoWait()
}

func EnvPath(path string) Option {
	return func(r *Runtime) {
		r.opts.Set(optEnvPath, path)
	}
}

// Deprecated: use EnvPath
func EnvPathOption(path string) Option {
	return EnvPath(path)
}

func WithConfig(cfg interface{}) Option {
	return func(r *Runtime) {
		r.cfg = cfg
	}
}

// Deprecated: use WithConfig
func ConfigOption(cfg interface{}) Option {
	return WithConfig(cfg)
}
