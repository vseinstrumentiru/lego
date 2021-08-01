package newrelic

import (
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/log"
)

func NewHandler(app *newrelic.Application) log.EntryHandler {
	return &handler{app: app}
}

type handler struct {
	app *newrelic.Application
}

func (h *handler) Name() string {
	return "newrelic"
}

func (h *handler) LevelEnabled(level logur.Level) bool {
	return log.LevelEnabled(level, logur.Error)
}

func (h *handler) Handle(msg log.Entry) {
	if err, ok := msg.(error); ok {
		txn := h.app.StartTransaction("error")
		txn.NoticeError(nrpkgerrors.Wrap(err))
		txn.End()
	}
}

func (h *handler) StopPropagation() bool {
	return false
}
