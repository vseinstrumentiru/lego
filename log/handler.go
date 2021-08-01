package log

import (
	"logur.dev/logur"
)

type Propagation struct {
	Stop     bool
	Priority int
}

type filterHandler struct {
	matcher EntryMatcher
	handler EntryHandler
}

func (h filterHandler) Name() string {
	return "filter"
}

func (h filterHandler) StopPropagation() bool {
	return h.handler.StopPropagation()
}

func (h filterHandler) LevelEnabled(level logur.Level) bool {
	return h.handler.LevelEnabled(level)
}

func (h filterHandler) Handle(notification Entry) {
	if h.matcher(notification) {
		return
	}

	h.handler.Handle(notification)
}

type errFilterHandler struct {
	matcher EntryErrMatcher
	handler EntryHandler
}

func (h errFilterHandler) Name() string {
	return "filter.error"
}

func (h errFilterHandler) StopPropagation() bool {
	return h.handler.StopPropagation()
}

func (h errFilterHandler) LevelEnabled(level logur.Level) bool {
	return h.handler.LevelEnabled(level)
}

func (h errFilterHandler) Handle(notification Entry) {
	if err, ok := notification.(error); ok && h.matcher(err) {
		return
	}

	h.handler.Handle(notification)
}

type compositeHandler []EntryHandler

func (h compositeHandler) Name() string {
	return "composite"
}

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

func (h compositeHandler) Handle(notification Entry) {
	for _, handler := range h {
		if handler.LevelEnabled(notification.Level()) {
			handler.Handle(notification)

			if handler.StopPropagation() {
				return
			}
		}
	}
}

func WithFilter(handler EntryHandler, matcher EntryMatcher) EntryHandler {
	return filterHandler{
		matcher: matcher,
		handler: handler,
	}
}

func WithErrFilter(handler EntryHandler, matcher EntryErrMatcher) EntryHandler {
	return errFilterHandler{
		matcher: matcher,
		handler: handler,
	}
}

func AppendHandler(parent EntryHandler, add EntryHandler) EntryHandler {
	if c, ok := parent.(compositeHandler); ok {
		c = append(c, add)
		parent = c
	} else {
		parent = compositeHandler{parent, add}
	}

	return parent
}
