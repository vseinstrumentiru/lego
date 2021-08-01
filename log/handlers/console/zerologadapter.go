package console

import (
	"context"

	"github.com/rs/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/log"
)

// LoggerAdapter is a Logur adapter for zerolog.
type LoggerAdapter struct {
	logger zerolog.Logger
}

// NewAdapter returns a new Logur logger.
func NewAdapter(logger zerolog.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: logger,
	}
}

// Trace implements the Logur Logger interface.
func (l *LoggerAdapter) Trace(msg string, fields ...map[string]interface{}) {
	sendEvent(l.logger.Trace(), msg, fields...)
}

// Debug implements the Logur Logger interface.
func (l *LoggerAdapter) Debug(msg string, fields ...map[string]interface{}) {
	sendEvent(l.logger.Debug(), msg, fields...)
}

// Info implements the Logur Logger interface.
func (l *LoggerAdapter) Info(msg string, fields ...map[string]interface{}) {
	sendEvent(l.logger.Info(), msg, fields...)
}

// Warn implements the Logur Logger interface.
func (l *LoggerAdapter) Warn(msg string, fields ...map[string]interface{}) {
	sendEvent(l.logger.Warn(), msg, fields...)
}

// Error implements the Logur Logger interface.
func (l *LoggerAdapter) Error(msg string, fields ...map[string]interface{}) {
	sendEvent(l.logger.Error(), msg, fields...)
}

func (l *LoggerAdapter) TraceContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Trace(msg, fields...)
}

func (l *LoggerAdapter) DebugContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Debug(msg, fields...)
}

func (l *LoggerAdapter) InfoContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Info(msg, fields...)
}

func (l *LoggerAdapter) WarnContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Warn(msg, fields...)
}

func (l *LoggerAdapter) ErrorContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Error(msg, fields...)
}

func (l *LoggerAdapter) LevelEnabled(level logur.Level) bool {
	return log.LevelEnabled(level, ToLogurLevel(l.logger.GetLevel()))
}

func ToLogurLevel(level zerolog.Level) logur.Level {
	switch level {
	case zerolog.TraceLevel:
		return logur.Trace
	case zerolog.DebugLevel:
		return logur.Debug
	case zerolog.InfoLevel:
		return logur.Info
	case zerolog.WarnLevel:
		return logur.Warn
	case zerolog.Disabled, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.NoLevel, zerolog.PanicLevel:
		return logur.Error
	}

	return logur.Error
}

func ToZerologLevel(level logur.Level) zerolog.Level {
	switch level {
	case logur.Trace:
		return zerolog.TraceLevel
	case logur.Debug:
		return zerolog.DebugLevel
	case logur.Info:
		return zerolog.InfoLevel
	case logur.Warn:
		return zerolog.WarnLevel
	case logur.Error:
		return zerolog.ErrorLevel
	default:
		return zerolog.ErrorLevel
	}
}

func sendEvent(event *zerolog.Event, msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		event.Fields(fields[0])
	}

	event.Msg(msg)
}
