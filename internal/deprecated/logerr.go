package deprecated

import (
	"context"
	"fmt"

	"emperror.dev/emperror"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

// Deprecated: use multilog.Logger
type LogErr interface {
	// ===============
	// Logger
	// ===============
	Log() logur.LoggerFacade
	// Trace logs a Trace event.
	//
	// Even more fine-grained information than Debug events.
	// Loggers not supporting this level should fall back to Debug.
	Trace(msg string, fields ...map[string]interface{})
	// Debug logs a Debug event.
	//
	// A verbose series of information events.
	// They are useful when debugging the system.
	Debug(msg string, fields ...map[string]interface{})
	// Info logs an Info event.
	//
	// General information about what's happening inside the system.
	Info(msg string, fields ...map[string]interface{})
	// Warn logs a Warn(ing) event.
	//
	// Non-critical events that should be looked at.
	Warn(msg string, fields ...map[string]interface{})
	// Error logs an Error event.
	//
	// Critical events that require immediate attention.
	// Loggers commonly provide Fatal and Panic levels above Error level,
	// but exiting and panicing is out of scope for a logging library.
	Error(msg string, fields ...map[string]interface{})
	// TraceContext logs a Trace event.
	//
	// Even more fine-grained information than Debug events.
	// Loggers not supporting this level should fall back to Debug.
	TraceContext(ctx context.Context, msg string, fields ...map[string]interface{})
	// DebugContext logs a Debug event.
	//
	// A verbose series of information events.
	// They are useful when debugging the system.
	DebugContext(ctx context.Context, msg string, fields ...map[string]interface{})
	// InfoContext logs an Info event.
	//
	// General information about what's happening inside the system.
	InfoContext(ctx context.Context, msg string, fields ...map[string]interface{})
	// WarnContext logs a Warn(ing) event.
	//
	// Non-critical events that should be looked at.
	WarnContext(ctx context.Context, msg string, fields ...map[string]interface{})
	// ErrorContext logs an Error event.
	//
	// Critical events that require immediate attention.
	// Loggers commonly provide Fatal and Panic levels above Error level,
	// but exiting and panicing is out of scope for a logging library.
	ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{})
	// WithFields annotates a log with key-value pairs.
	WithFields(fields map[string]interface{}) LogErr
	// WithContext annotates a log with a contexttool.
	WithContext(ctx context.Context) LogErr
	// ===============
	// Error Handler
	// ===============
	Handler() emperror.ErrorHandlerFacade

	Handle(err error)

	HandleContext(ctx context.Context, err error)

	WithFilter(matcher multilog.EntryErrMatcher) LogErr

	WithDetails(details ...interface{}) LogErr
}

var _ LogErr = &multilogWrapper{}
var _ LogErrF = &multilogWrapper{}

// Deprecated: use multilog.Logger
func NewLogErr(logErr LogErr) LogErrF {
	if w, ok := logErr.(*multilogWrapper); ok {
		return w
	}

	return nil
}

type multilogWrapper struct {
	multilog.Logger
}

func (m *multilogWrapper) Errorf(format string, v ...interface{}) {
	m.Error(fmt.Sprintf(format, v...))
}

func (m *multilogWrapper) Warnf(format string, v ...interface{}) {
	m.Warn(fmt.Sprintf(format, v...))
}

func (m *multilogWrapper) Debugf(format string, v ...interface{}) {
	m.Debug(fmt.Sprintf(format, v...))
}

func (m *multilogWrapper) Log() logur.LoggerFacade {
	return m
}

func (m *multilogWrapper) WithFields(fields map[string]interface{}) LogErr {
	return &multilogWrapper{Logger: m.Logger.WithFields(fields)}
}

func (m *multilogWrapper) WithContext(ctx context.Context) LogErr {
	return &multilogWrapper{Logger: m.Logger.WithContext(ctx)}
}

func (m *multilogWrapper) Handler() emperror.ErrorHandlerFacade {
	return m
}

func (m *multilogWrapper) WithFilter(matcher multilog.EntryErrMatcher) LogErr {
	return &multilogWrapper{Logger: m.Logger.WithErrFilter(matcher)}
}

func (m *multilogWrapper) WithDetails(details ...interface{}) LogErr {
	return m
}

// Deprecated: use multilog.Logger
type LogErrF interface {
	LogErr
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}
