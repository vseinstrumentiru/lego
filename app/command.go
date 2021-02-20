package app

import (
	"os"

	"emperror.dev/emperror"
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	di "github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
	"github.com/vseinstrumentiru/lego/v2/version"
)

func rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "lego",
		Hidden: true,
	}
}

func showVersion(e env.Env, c di.ChainContainer) {
	e.OnFlag("version", func(bool) {
		c.Execute(func(ver *version.Info) { ver.Print() })
		os.Exit(0)
	})
}

type cmdArgs struct {
	dig.In
	Root     *cobra.Command
	Children []*cobra.Command `group:"cmd"`
}

func printDIGraph(e env.Env, c di.Container, cfg *config.Application) {
	if cfg.DebugMode {
		e.OnFlag("di", func(bool) {
			emperror.Panic(c.Visualize(os.Stdout))
			os.Exit(0)
		})
	}
}

func command(r *runtime) {
	r.container.Execute(func(in cmdArgs) error {
		in.Root.AddCommand(in.Children...)
		cmds := in.Root.Commands()
		if len(cmds) > 0 {
			return in.Root.Execute()
		}

		return nil
	})
}
