package multilogprovider

import (
	"context"

	"github.com/rs/zerolog"
	"logur.dev/logur"
)

// Logger is a Logur adapter for zerolog.
type LoggerAdapter struct {
	logger zerolog.Logger
}

// New returns a new Logur logger.
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
	return int(level) >= int(l.logger.GetLevel()+1)
}

func sendEvent(event *zerolog.Event, msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		event.Fields(fields[0])
	}

	event.Msg(msg)
}
