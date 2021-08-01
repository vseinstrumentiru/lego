package log

import (
	"context"

	"emperror.dev/emperror"
	"logur.dev/logur"
)

type Level = logur.Level

type Entry interface {
	Level() Level
	Message() string
	Fields() map[string]interface{}
	Details() []interface{}

	WithDetails(details ...interface{}) Entry
	WithFields(fields map[string]interface{}) Entry
}

type Logger interface {
	logur.LoggerFacade
	emperror.ErrorHandlerFacade
	Notify(notification interface{})
	WithLevel(level logur.Level) Logger
	WithHandler(handler EntryHandler) Logger
	WithContext(ctx context.Context) Logger
	WithFilter(matcher EntryMatcher) Logger
	WithErrFilter(matcher EntryErrMatcher) Logger
	WithFields(fields map[string]interface{}) Logger
}

type EntryHandler interface {
	logur.LevelEnabler
	Name() string
	Handle(msg Entry)
	StopPropagation() bool
}

type ContextExtractor func(ctx context.Context) map[string]interface{}

type (
	EntryMatcher    func(msg Entry) bool
	EntryErrMatcher func(err error) bool
)
