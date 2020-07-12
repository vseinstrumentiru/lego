package notify

import (
	"logur.dev/logur"
)

type Level = logur.Level

type LoggerFacade interface {
	logur.LoggerFacade
	logur.LevelEnabler
}

type CustomLoggerFacade interface {
	Print(v ...interface{})
}

type NotificationHandler interface {
	LoggerFacade
	CustomLoggerFacade
}
