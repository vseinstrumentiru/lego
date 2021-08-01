package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/log"
)

type logger struct {
	log.Logger
}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var lvl logur.Level
	switch level {
	case pgx.LogLevelTrace, pgx.LogLevelDebug, pgx.LogLevelInfo:
		lvl = logur.Debug
	case pgx.LogLevelWarn:
		lvl = logur.Warn
	case pgx.LogLevelError:
		lvl = logur.Error
	default:
		return
	}

	l.Logger.WithContext(ctx).Notify(log.NewEntry(lvl, msg, data))
}
