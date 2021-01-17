package app

import (
	"os"
	"os/signal"
	"syscall"

	"emperror.dev/errors/match"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
)

func serve(r *runtime) {
	log := r.log
	log.Trace("starting pipeline")
	var pipeline *run.Group
	var upg *tableflip.Upgrader
	r.container.Execute(func(p *run.Group, u *tableflip.Upgrader) {
		pipeline = p
		upg = u
	})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			log.Info("graceful reloading")
			log.Notify(upg.Upgrade())
		}
	}()

	// running application
	if err := pipeline.Run(); err != nil {
		log.WithErrFilter(match.As(&run.SignalError{}).MatchError).Notify(err)
	}
}
