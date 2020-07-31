package notify

import (
	"context"
	"emperror.dev/errors"
	"logur.dev/logur"
)

type Option func(c notificator) notificator

func WithHandler(handler i.NotificationHandler) Option {
	return func(c notificator) notificator {
		c.handler = AppendHandler(c.handler, handler)

		return c
	}
}

func New(level i.Level, options ...Option) i.Notificator {
	n := notificator{
		level:   level,
		handler: compositeHandler{},
	}

	for i := 0; i < len(options); i++ {
		n = options[i](n)
	}

	return n
}

type notificator struct {
	level     logur.Level
	handler   i.NotificationHandler
	fields    map[string]interface{}
	extractor i.ContextExtractor
}

func (c notificator) WithFields(fields map[string]interface{}) i.Notificator {
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

	return notificator{
		level:     c.level,
		handler:   c.handler,
		fields:    fields,
		extractor: c.extractor,
	}
}

func (c notificator) LevelEnabled(level logur.Level) bool {
	return c.level >= level
}

func (c notificator) Notify(notification interface{}) {
	var n i.Notification

	switch t := notification.(type) {
	case i.Notification:
		n = t
	case error:
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

func (c notificator) WithFilter(matcher i.NotificationMatcher) i.Notificator {
	return &notificator{
		level:     c.level,
		handler:   WithFilter(c.handler, matcher),
		extractor: c.extractor,
	}
}

func (c notificator) WithContext(ctx context.Context) i.Notificator {
	if c.extractor == nil {
		return c
	}
	return c.WithFields(c.extractor(ctx))
}

func (c notificator) notify(level logur.Level, msg string, fields ...map[string]interface{}) {
	if c.level < level {
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

func (c notificator) Trace(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Trace, msg, fields...)
}

func (c notificator) Debug(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Debug, msg, fields...)
}

func (c notificator) Info(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Info, msg, fields...)
}

func (c notificator) Warn(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Warn, msg, fields...)
}

func (c notificator) Error(msg string, fields ...map[string]interface{}) {
	c.notify(logur.Error, msg, fields...)
}

func (c notificator) notifyContext(ctx context.Context, level logur.Level, msg string, fields ...map[string]interface{}) {
	if c.extractor != nil {
		fields = append(fields, c.extractor(ctx))
	}

	c.notify(level, msg, fields...)
}

func (c notificator) TraceContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Trace, msg, fields...)
}

func (c notificator) DebugContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Debug, msg, fields...)
}

func (c notificator) InfoContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Info, msg, fields...)
}

func (c notificator) WarnContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Warn, msg, fields...)
}

func (c notificator) ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	c.notifyContext(ctx, logur.Error, msg, fields...)
}
