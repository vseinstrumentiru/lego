package provide

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
)

const (
	defaultEnvPath = "app"
)

type envArgs struct {
	dig.In
	Runtime config.Runtime
	Config  env.Config `optional:"true"`
	Cmd     *cobra.Command
}

func Env(in envArgs) (env.Env, config.Env) {
	path := defaultEnvPath

	in.Runtime.On(config.OptEnvPath, func(newPath string) {
		path = newPath
	})

	var instance env.Env

	if in.Config == nil {
		instance = env.NewNoConfigEnv(in.Cmd.PersistentFlags(), path)
	} else {
		instance = env.NewConfigEnv(in.Cmd.PersistentFlags(), path)
	}

	instance.SetFlag("version", false, "show version")
	if in.Runtime.Is(config.OptLocalDebug) {
		instance.SetFlag("di", false, "show dependency graph (debug mode only)")
	}

	return instance, instance
}
