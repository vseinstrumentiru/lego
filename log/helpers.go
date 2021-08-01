package log

import (
	"log"

	"logur.dev/logur"
)

// NewErrorStandardLogger returns a new standard log logging on error level.
func NewErrorStandardLogger(logger Logger) *log.Logger {
	return logur.NewErrorStandardLogger(logger, "", 0)
}

// SetStandardLogger sets the global log's output to a custom log instance.
func SetStandardLogger(logger Logger) {
	log.SetOutput(logur.NewWriter(logger))
}

// WithHandler option adds handler to logger
func WithHandler(handler EntryHandler) Option {
	return func(c *logger) *logger {
		c.handler = AppendHandler(c.handler, handler)

		return c
	}
}

// NoopLogger returns empty logger
func NoopLogger() Logger {
	return New(logur.Error)
}

func LevelEnabled(entryLevel logur.Level, handlerLevel logur.Level) bool {
	return entryLevel >= handlerLevel
}
