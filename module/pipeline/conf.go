package pipeline

import (
	"context"
	"syscall"

	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"go.uber.org/dig"
)

type Args struct {
	dig.In
	Pipeline *run.Group
	Upg      *tableflip.Upgrader `optional:"true"`
}

func Configure(in Args) {
	ctx := context.Background()
	in.Pipeline.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	if in.Upg != nil {
		in.Pipeline.Add(appkitrun.GracefulRestart(ctx, in.Upg))
	}
}
