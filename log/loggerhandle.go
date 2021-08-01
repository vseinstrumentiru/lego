package log

import (
	"emperror.dev/errors"
	"logur.dev/logur"
)

func (l *logger) Notify(notification interface{}) {
	var n Entry

	switch t := notification.(type) {
	case Entry:
		if t == nil {
			return
		}
		n = t.WithFields(l.fields)
	case error:
		if t == nil {
			return
		}

		n = NewErrEntry(errors.WithStackDepthIf(t, 1), l.fields)
	case string:
		n = NewEntry(logur.Info, t, l.fields)
	default:
		return
	}

	if l.LevelEnabled(n.Level()) {
		l.handler.Handle(n)
	}
}

func (l *logger) Handle(err error) {
	l.Notify(err)
}

func (l *logger) Trace(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Trace, msg, fields...)
}

func (l *logger) Debug(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Debug, msg, fields...)
}

func (l *logger) Info(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Info, msg, fields...)
}

func (l *logger) Warn(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Warn, msg, fields...)
}

func (l *logger) Error(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Error, msg, fields...)
}

func (l *logger) notify(level logur.Level, msg string, fields ...map[string]interface{}) {
	if !l.LevelEnabled(level) {
		return
	}

	_fields := make(map[string]interface{})

	if len(l.fields) != 0 {
		for key, val := range l.fields {
			_fields[key] = val
		}
	}

	for i := 0; i < len(fields); i++ {
		for key, val := range fields[i] {
			_fields[key] = val
		}
	}

	n := NewEntry(level, msg, _fields)

	l.Notify(n)
}
