package execute

import (
	"os"

	"emperror.dev/emperror"
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	. "github.com/vseinstrumentiru/lego/v2/config"
	di "github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/internal/env"
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/opencensusexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/prometheus"
	"github.com/vseinstrumentiru/lego/v2/version"
)

func All(r Runtime) (exec []interface{}) {
	var res []interface{}
	if r.Is(OptLocalDebug) {
		res = append(res, PrintDIGraph)
	}
	if r.Is(OptWithoutProviders) {
		return res
	}

	if r.Is(ServerMode) {
		res = append(res, Pipeline, prometheus.Configure)
	}

	res = append(res,
		jaegerexporter.Configure,
		opencensusexporter.Configure,
		newrelicexporter.Configure,
		metrics.ConfigureTrace,
		metrics.ConfigureStats,
	)

	return res
}

type cmdArgs struct {
	dig.In
	Root     *cobra.Command
	Children []*cobra.Command `group:"cmd"`
}

func Version(e env.Env, c di.ChainContainer) {
	e.OnFlag("version", func(bool) {
		c.Execute(func(ver *version.Info) { ver.Print() })
		os.Exit(0)
	})
}

func RunCommands(in cmdArgs) error {
	in.Root.AddCommand(in.Children...)
	cmds := in.Root.Commands()
	if len(cmds) > 0 {
		return in.Root.Execute()
	}

	return nil
}

func PrintDIGraph(e env.Env, c di.Container, cfg *Application) {
	if cfg.DebugMode {
		e.OnFlag("di", func(bool) {
			emperror.Panic(c.Visualize(os.Stdout))
			os.Exit(0)
		})
	}
}
