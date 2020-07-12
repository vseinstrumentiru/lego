package server

const (
	flagNoWait = "flag_nowait"
	optEnvPath = "opt_env_path"
)

func NoWaitOption() Option {
	return func(r *Runtime) {
		r.opts.Set(flagNoWait, true)
	}
}

func EnvPathOption(path string) Option {
	return func(r *Runtime) {
		r.opts.Set(optEnvPath, path)
	}
}

func ConfigOption(cfg interface{}) Option {
	return func(r *Runtime) {
		r.cfg = cfg
	}
}
