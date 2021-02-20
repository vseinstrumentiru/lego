package env

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
)

const (
	defaultEnvPath = "app"
)

type envArgs struct {
	dig.In
	Runtime config.Runtime
	Config  Config `optional:"true"`
	Cmd     *cobra.Command
}

func Provide(in envArgs) (Env, config.Env) {
	path := defaultEnvPath

	in.Runtime.On(config.OptEnvPath, func(newPath string) {
		path = newPath
	})

	var instance Env

	if in.Config == nil {
		instance = NewNoConfigEnv(in.Cmd.PersistentFlags(), path)
	} else {
		instance = NewConfigEnv(in.Cmd.PersistentFlags(), path)
	}

	instance.SetFlag("version", false, "show version")
	if in.Runtime.Is(config.OptLocalDebug) {
		instance.SetFlag("di", false, "show dependency graph (debug mode only)")
	}

	return instance, instance
}
