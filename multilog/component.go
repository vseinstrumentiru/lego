package multilog

import (
	"context"

	"emperror.dev/errors"
	"logur.dev/logur"
)

type Option func(c multilog) multilog

func WithHandler(handler EntryHandler) Option {
	return func(c multilog) multilog {
		c.handler = AppendHandler(c.handler, handler)

		return c
	}
}

func NoopLogger() Logger {
	return New(logur.Error)
}

func New(level Level, options ...Option) Logger {
	n := multilog{
		level:   level,
		handler: compositeHandler{},
	}

	for i := 0; i < len(options); i++ {
		n = options[i](n)
	}

	return n
}

type multilog struct {
	level     logur.Level
	handler   EntryHandler
	fields    map[string]interface{}
	extractor ContextExtractor
}

func (l multilog) WithHandler(handler EntryHandler) {
	l.handler = AppendHandler(l.handler, handler)
}

func (l multilog) HandleContext(ctx context.Context, err error) {
	l.WithContext(ctx).Notify(err)
}

func (l multilog) Handle(err error) {
	l.Notify(err)
}

func (l multilog) WithFields(fields map[string]interface{}) Logger {
	if len(l.fields) > 0 {
		_fields := make(map[string]interface{}, len(l.fields)+len(fields))

		for key, value := range l.fields {
			_fields[key] = value
		}

		for key, value := range fields {
			_fields[key] = value
		}

		fields = _fields
	}

	return multilog{
		level:     l.level,
		handler:   l.handler,
		fields:    fields,
		extractor: l.extractor,
	}
}

func (l multilog) LevelEnabled(level logur.Level) bool {
	return l.level >= level
}

func (l multilog) Notify(notification interface{}) {
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
		n = errorNotification{
			err:    errors.WithStackDepthIf(t, 1),
			fields: l.fields,
		}
	case string:
		n = NewNotification(logur.Info, t, l.fields)
	default:
		return
	}

	l.handler.Handle(n)
}

func (l multilog) WithFilter(matcher EntryMatcher) Logger {
	return &multilog{
		level:     l.level,
		handler:   WithFilter(l.handler, matcher),
		extractor: l.extractor,
	}
}

func (l multilog) WithErrFilter(matcher EntryErrMatcher) Logger {
	return &multilog{
		level:     l.level,
		handler:   WithErrFilter(l.handler, matcher),
		extractor: l.extractor,
	}
}

func (l multilog) WithContext(ctx context.Context) Logger {
	if l.extractor == nil {
		return l
	}
	return l.WithFields(l.extractor(ctx))
}

func (l multilog) notify(level logur.Level, msg string, fields ...map[string]interface{}) {
	if l.level > level {
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

	n := NewNotification(level, msg, _fields)

	l.Notify(n)
}

func (l multilog) Trace(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Trace, msg, fields...)
}

func (l multilog) Debug(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Debug, msg, fields...)
}

func (l multilog) Info(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Info, msg, fields...)
}

func (l multilog) Warn(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Warn, msg, fields...)
}

func (l multilog) Error(msg string, fields ...map[string]interface{}) {
	l.notify(logur.Error, msg, fields...)
}

func (l multilog) notifyContext(ctx context.Context, level logur.Level, msg string, fields ...map[string]interface{}) {
	if l.extractor != nil {
		fields = append(fields, l.extractor(ctx))
	}

	l.notify(level, msg, fields...)
}

func (l multilog) TraceContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Trace, msg, fields...)
}

func (l multilog) DebugContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Debug, msg, fields...)
}

func (l multilog) InfoContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Info, msg, fields...)
}

func (l multilog) WarnContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Warn, msg, fields...)
}

func (l multilog) ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Error, msg, fields...)
}
