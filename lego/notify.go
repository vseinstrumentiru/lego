package lego

import (
	"context"
	"logur.dev/logur"
)

type Level = logur.Level

type Notification interface {
	Level() Level
	Message() string
	Fields() map[string]interface{}
	Details() []interface{}

	WithDetails(details ...interface{}) Notification
}

type Notificator interface {
	logur.LoggerFacade
	Notify(notification interface{})
	WithContext(ctx context.Context) Notificator
	WithFilter(matcher NotificationMatcher) Notificator
	WithFields(fields map[string]interface{}) Notificator
}

type NotificationHandler interface {
	logur.LevelEnabler
	Handle(msg Notification)
	StopPropagation() bool
}

type ContextExtractor func(ctx context.Context) map[string]interface{}

type NotificationMatcher func(msg Notification) bool
