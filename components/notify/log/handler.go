package log

import (
	"logur.dev/logur"
)

func Handler(logger logur.Logger, stop bool) *handler {
	return &handler{
		handler: logger,
		stop:    stop,
	}
}

type handler struct {
	handler logur.Logger
	stop    bool
}

func (h *handler) LevelEnabled(level logur.Level) bool {
	if en, ok := h.handler.(logur.LevelEnabler); ok {
		return en.LevelEnabled(level)
	}

	return true
}

func (h *handler) Handle(msg i.Notification) {
	switch msg.Level() {
	case logur.Trace:
		h.handler.Trace(msg.Message(), msg.Fields())
	case logur.Debug:
		h.handler.Debug(msg.Message(), msg.Fields())
	case logur.Info:
		h.handler.Info(msg.Message(), msg.Fields())
	case logur.Warn:
		h.handler.Warn(msg.Message(), msg.Fields())
	case logur.Error:
		h.handler.Error(msg.Message(), msg.Fields())
	}
}

func (h *handler) StopPropagation() bool {
	return h.stop
}
