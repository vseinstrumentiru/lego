package log

import (
	"context"

	"logur.dev/logur"
)

func (l *logger) HandleContext(ctx context.Context, err error) {
	l.WithContext(ctx).Notify(err)
}

func (l *logger) TraceContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Trace, msg, fields...)
}

func (l *logger) DebugContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Debug, msg, fields...)
}

func (l *logger) InfoContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Info, msg, fields...)
}

func (l *logger) WarnContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Warn, msg, fields...)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{}) {
	l.notifyContext(ctx, logur.Error, msg, fields...)
}

func (l *logger) notifyContext(ctx context.Context, level logur.Level, msg string, fields ...map[string]interface{}) {
	if l.extractor != nil {
		fields = append(fields, l.extractor(ctx))
	}

	l.notify(level, msg, fields...)
}
