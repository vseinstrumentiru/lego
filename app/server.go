package app

import (
	"emperror.dev/errors/match"
	"github.com/oklog/run"
)

func serve(r *runtime) {
	log := r.log
	log.Trace("starting pipeline")
	var pipeline *run.Group
	r.container.Execute(func(p *run.Group) {
		pipeline = p
	})

	// running application
	if err := pipeline.Run(); err != nil {
		log.WithErrFilter(match.As(&run.SignalError{}).MatchError).Notify(err)
	}
}
