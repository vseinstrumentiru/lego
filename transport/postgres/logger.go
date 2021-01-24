package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type logger struct {
	multilog.Logger
}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var lvl logur.Level
	switch level {
	case pgx.LogLevelTrace:
		lvl = logur.Trace
	case pgx.LogLevelDebug:
		lvl = logur.Debug
	case pgx.LogLevelInfo:
		lvl = logur.Info
	case pgx.LogLevelWarn:
		lvl = logur.Warn
	case pgx.LogLevelError:
		lvl = logur.Error
	default:
		return
	}

	l.Logger.WithContext(ctx).Notify(multilog.NewNotification(lvl, msg, data))
}
