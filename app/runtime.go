package app

import (
	"emperror.dev/emperror"
	"github.com/spf13/cobra"

	"github.com/vseinstrumentiru/lego/v2/common/cast"
	"github.com/vseinstrumentiru/lego/v2/common/set"
	"github.com/vseinstrumentiru/lego/v2/config"
	di "github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/log/handlers/console"
	"github.com/vseinstrumentiru/lego/v2/log/handlers/sentry"
	"github.com/vseinstrumentiru/lego/v2/log/logger"
	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type runner func(r *runtime)

func NewRuntime(opts ...Option) config.Runtime {
	parent := di.New()
	container := di.NewChain(parent)
	rt := &runtime{
		opts:      cast.NewCastableRWSet(set.NewSimple()),
		container: container,
		exec:      command,
	}

	for _, opt := range opts {
		opt(rt)
	}

	container.Register(func() config.Runtime { return rt })
	container.Register(func() di.ChainContainer { return container })
	container.Register(func() di.Container { return parent })
	container.Make(rt)

	return rt
}

type runtime struct {
	container      di.ChainContainer
	configurations []interface{}
	log            log.Logger
	opts           cast.CastableRWSet
	exec           runner
}

func (r *runtime) Providers() []interface{} {
	return []interface{}{
		rootCommand,
		version.New,
		shutdown.NewCloseGroup,
		logger.Provide,
		sentry.Provide,
		console.Provide,
		env.Provide,
	}
}

func (r *runtime) Configurations() []interface{} {
	exec := []interface{}{
		r.configureEnv,
		printDIGraph,
		r.configureVersion,
		showVersion,
		r.configureLogger,
	}

	return append(exec, r.configurations...)
}

func (r *runtime) configureEnv(e env.Env) {
	r.container.Make(e)
}

func (r *runtime) configureVersion(ver *version.Info, cfg *config.Application, cmd *cobra.Command) {
	ver.DataCenter = cfg.DataCenter
	if !config.IsUndefined(cfg.Name) {
		cmd.Use = cfg.Name
	}

	if r.Is(config.OptLocalDebug) {
		cfg.LocalMode = true
		cfg.DebugMode = true
	}
}

func (r *runtime) configureLogger(logger log.Logger) {
	r.log = logger.WithFields(map[string]interface{}{"component": "runtime"})
	log.SetStandardLogger(r.log.WithFields(map[string]interface{}{"type": "standard"}))
	r.log.Trace("starting application")
}

func (r *runtime) On(key string, callback interface{}) (ok bool) {
	return cast.OnCheck(r.opts, key, callback)
}

func (r *runtime) Is(key string) bool {
	ok, err := r.opts.GetBool(key)

	return ok && err == nil
}

func (r *runtime) Not(key string) bool {
	ok, err := r.opts.GetBool(key)

	return !ok || err != nil
}

func (r *runtime) Run(apps ...interface{}) {
	defer emperror.HandleRecover(r.log.WithFields(map[string]interface{}{"type": "panic"}))
	// pipeline configuration
	var closer *shutdown.CloseGroup
	log := r.log
	r.container.Execute(func(c *shutdown.CloseGroup) { closer = c })
	defer func() {
		err := closer.Close()
		log.Notify(err)
	}()

	// constructing application
	log.Trace("constructing application")
	for i := 0; i < len(apps); i++ {
		r.container.Make(apps[i])
	}
	log.Trace("constructing application", map[string]interface{}{"status": "completed"})

	r.exec(r)

	log.Trace("application stopped")
}
