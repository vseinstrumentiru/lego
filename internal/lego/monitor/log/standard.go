package log

import (
	"log"
	"logur.dev/logur"
)

// NewErrorStandardLogger returns a new standard log logging on error level.
func NewErrorStandardLogger(logger logur.Logger) *log.Logger {
	return logur.NewErrorStandardLogger(logger, "", 0)
}

// SetStandardLogger sets the global log's output to a custom log instance.
func SetStandardLogger(logger logur.Logger) {
	log.SetOutput(logur.NewLevelWriter(logger, logur.Info))
}
