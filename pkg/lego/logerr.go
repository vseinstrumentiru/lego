package lego

import (
	"context"
	"emperror.dev/emperror"
	"fmt"
	"logur.dev/logur"
)

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

	WithFilter(matcher ErrorMatcher) LogErr

	WithDetails(details ...interface{}) LogErr
}

type ErrorMatcher func(err error) bool

type ContextExtractor func(ctx context.Context) map[string]interface{}

type logErr struct {
	logur.LoggerFacade
	emperror.ErrorHandlerFacade
	extractor ContextExtractor
}

func NewLogErr(logger logur.LoggerFacade, handler emperror.ErrorHandlerFacade) LogErr {
	return &logErr{
		LoggerFacade:       logger,
		ErrorHandlerFacade: handler,
	}
}

// NewContextAwareLogger returns a new Logger instance that can extract information from a contexttool.
func NewContextAwareLogErr(logger logur.LoggerFacade, handler emperror.ErrorHandlerFacade, extractor ContextExtractor) LogErr {
	return &logErr{
		LoggerFacade:       logur.WithContextExtractor(logger, logur.ContextExtractor(extractor)),
		ErrorHandlerFacade: emperror.WithContextExtractor(handler, emperror.ContextExtractor(extractor)),
		extractor:          extractor,
	}
}

func (l *logErr) Log() logur.LoggerFacade {
	return l.LoggerFacade
}

func (l *logErr) Handler() emperror.ErrorHandlerFacade {
	return l.ErrorHandlerFacade
}

func (l *logErr) WithFields(fields map[string]interface{}) LogErr {
	return &logErr{
		LoggerFacade:       logur.WithFields(l.LoggerFacade, fields),
		ErrorHandlerFacade: l.ErrorHandlerFacade,
		extractor:          l.extractor,
	}
}

// WithContext annotates a log with a contexttool.
func (l *logErr) WithContext(ctx context.Context) LogErr {
	if l.extractor == nil {
		return l
	}

	return l.WithFields(l.extractor(ctx))
}

// WithContext annotates a log with a contexttool.
func (l *logErr) WithFilter(matcher ErrorMatcher) LogErr {
	return &logErr{
		LoggerFacade:       l.LoggerFacade,
		ErrorHandlerFacade: emperror.WithFilter(l.ErrorHandlerFacade, emperror.ErrorMatcher(matcher)),
		extractor:          l.extractor,
	}
}

func (l *logErr) WithDetails(details ...interface{}) LogErr {
	return &logErr{
		LoggerFacade:       l.LoggerFacade,
		ErrorHandlerFacade: emperror.WithDetails(l.ErrorHandlerFacade, details...),
		extractor:          l.extractor,
	}
}

type LogErrF interface {
	LogErr
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type logErrF struct {
	LogErr
}

func NewLogErrF(logErr LogErr) *logErrF {
	return &logErrF{LogErr: logErr}
}

func (l *logErrF) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l *logErrF) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l *logErrF) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}
