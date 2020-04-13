package log

import (
	"logur.dev/logur"
)

func Provide(config Config, env string, appName string) (logger logur.LoggerFacade) {
	logger = New(config)

	// Provide some basic contexttool to all log lines
	logger = logur.WithFields(logger, map[string]interface{}{"environment": env, "application": appName})
	SetStandardLogger(logger)

	return
}
