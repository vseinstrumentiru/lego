package log

import (
	"logur.dev/logur"
)

type Option func(c *logger) *logger

func New(level Level, options ...Option) Logger {
	n := &logger{
		level:   level,
		handler: compositeHandler{},
	}

	for i := 0; i < len(options); i++ {
		n = options[i](n)
	}

	return n
}

type logger struct {
	level     logur.Level
	handler   EntryHandler
	fields    map[string]interface{}
	extractor ContextExtractor
}

func (l *logger) LevelEnabled(level logur.Level) bool {
	return level >= l.level
}
