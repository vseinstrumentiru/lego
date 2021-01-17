package execute

import (
	"context"
	"syscall"

	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
)

func Pipeline(pipeline *run.Group, upg *tableflip.Upgrader) {
	ctx := context.Background()
	pipeline.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	pipeline.Add(appkitrun.GracefulRestart(ctx, upg))
}
