package notify

import (
	"github.com/vseinstrumentiru/lego/lego"
	"logur.dev/logur"
)

type filterHandler struct {
	matcher lego.NotificationMatcher
	handler lego.NotificationHandler
}

type Propagation struct {
	Stop     bool
	Priority int
}

func (h filterHandler) StopPropagation() bool {
	return h.handler.StopPropagation()
}

func (h filterHandler) LevelEnabled(level logur.Level) bool {
	return h.handler.LevelEnabled(level)
}

func (h filterHandler) Handle(notification lego.Notification) {
	if h.matcher(notification) {
		return
	}

	h.handler.Handle(notification)
}

type compositeHandler []lego.NotificationHandler

func (h compositeHandler) StopPropagation() bool {
	return false
}

func (h compositeHandler) LevelEnabled(level logur.Level) bool {
	for _, handler := range h {
		if handler.LevelEnabled(level) {
			return true
		}
	}

	return false
}

func (h compositeHandler) Handle(notification lego.Notification) {
	for _, handler := range h {
		if handler.LevelEnabled(notification.Level()) {
			handler.Handle(notification)

			if handler.StopPropagation() {
				return
			}
		}
	}
}

func WithFilter(handler lego.NotificationHandler, matcher lego.NotificationMatcher) lego.NotificationHandler {
	return filterHandler{
		matcher: matcher,
		handler: handler,
	}
}

func AppendHandler(parent lego.NotificationHandler, add lego.NotificationHandler) lego.NotificationHandler {
	if c, ok := parent.(compositeHandler); ok {
		c = append(c, add)
	} else {
		parent = compositeHandler{parent, add}
	}

	return parent
}
