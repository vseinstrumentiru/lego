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

func (c multilog) HandleContext(ctx context.Context, err error) {
	c.WithContext(ctx).Notify(err)
}

func (c multilog) Handle(err error) {
	c.Notify(err)
}

func (c multilog) WithFields(fields map[string]interface{}) Logger {
	if len(c.fields) > 0 {
		_fields := make(map[string]interface{}, len(c.fields)+len(fields))

		for key, value := range c.fields {
			_fields[key] = value
		}

		for key, value := range fields {
			_fields[key] = value
		}

		fields = _fields
	}

	return multilog{
		level:     c.level,
		handler:   c.handler,
		fields:    fields,
		extractor: c.extractor,
	}
}

func (c multilog) LevelEnabled(level logur.Level) bool {
	return c.level >= level
}

func (c multilog) Notify(notification interface{}) {
	var n Entry

	switch t := notification.(type) {
	case Entry:
		n = t
	case error:
		if t == nil {
			return
		}
		n = errorNotification{
			err:    errors.WithStackDepthIf(t, 1),
			fields: c.fields,
		}
	case string:
		n = NewNotification(logur.Info, t, c.fields)
	default:
		return
	}

	c.handler.Handle(n)
}

func (c multilog) WithFilter(matcher EntryMatcher) Logger {
	return &multilog{
		level:     c.level,
		handler:   WithFilter(c.handler, matcher),
		extractor: c.extractor,
	}
}

func (c multilog) WithErrFilter(matcher EntryErrMatcher) Logger {
	return &multilog{
		level:     c.level,
		handler:   WithErrFilter(c.handler, matcher),
		extractor: c.extractor,
	}
}

func (c multilog) WithContext(ctx context.Context) Logger {
	if c.extractor == nil {
		return c
	}
	return c.WithFields(c.extractor(ctx))
}

func (c multilog) notify(level logur.Level, msg string, fields ...map[string]interface{}) {
	if c.level > level {
		return
	}

	_fields := make(map[string]interface{})

	if len(c.fields) != 0 {
		for key, val := range c.fields {
			_fields[key] = val
		}
	}

	for i := 0; i < len(fields); i++ {
		for key, val := range fields[i] {
			_fields[key] = val
		}
	}

	n := NewNotification(level, msg, _fields)

	c.Notify(n)
}

func (c multilog) Trace(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Trace, msg, fields...)
}

func (c multilog) Debug(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Debug, msg, fields...)
}

func (c multilog) Info(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Info, msg, fields...)
}

func (c multilog) Warn(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Warn, msg, fields...)
}

func (c multilog) Error(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Error, msg, fields...)
}

func (c multilog) notifyContext(ctx context.Context, level logur.Level, msg string, fields ...map[string]interface{}) {
	if c.extractor != nil {
		fields = append(fields, c.extractor(ctx))
	}

	c.notify(level, msg, fields...)
}

func (c multilog) TraceContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Trace, msg, fields...)
}

func (c multilog) DebugContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Debug, msg, fields...)
}

func (c multilog) InfoContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Info, msg, fields...)
}

func (c multilog) WarnContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Warn, msg, fields...)
}

func (c multilog) ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Error, msg, fields...)
}
