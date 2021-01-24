package newrelic

import (
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

func Handler(app *newrelic.Application) *handler {
	return &handler{app: app}
}

type handler struct {
	app *newrelic.Application
}

func (h *handler) LevelEnabled(level logur.Level) bool {
	return level >= logur.Error
}

func (h *handler) Handle(msg multilog.Entry) {
	if err, ok := msg.(error); ok {
		txn := h.app.StartTransaction("error")
		txn.NoticeError(nrpkgerrors.Wrap(err))
		txn.End()
	}
}

func (h *handler) StopPropagation() bool {
	return false
}
