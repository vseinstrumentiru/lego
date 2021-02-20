package uprader

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type Args struct {
	dig.In
	Upg *tableflip.Upgrader `optional:"true"`
	Log multilog.Logger
}

func Configure(in Args) error {
	if in.Upg == nil {
		return nil
	}
	log := in.Log
	upg := in.Upg

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			log.Info("graceful reloading")
			log.Notify(upg.Upgrade())
		}
	}()

	return nil
}
