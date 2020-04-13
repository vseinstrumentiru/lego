// Package log configures a new log for an application.
package log

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

// NewLogger creates a new log.
func New(config Config) logur.LoggerFacade {
	logger := logrus.New()
	logger.SetReportCaller(config.UseStack)

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:             config.NoColor,
		EnvironmentOverrideColors: true,
		CallerPrettyfier:          stackFn(config.SkipStack),
	})

	switch config.Format {
	case "logfmt":
		// Already the default

	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: stackFn(config.SkipStack),
		})
	}

	if level, err := logrus.ParseLevel(config.Level); err == nil {
		logger.SetLevel(level)
	}

	return logrusadapter.New(logger)
}

func stackFn(skip int) func(*runtime.Frame) (string, string) {
	return func(_ *runtime.Frame) (function string, file string) {
		pcs := callers(5)
		frames := runtime.CallersFrames(pcs[skip : skip+1])

		frame, more := frames.Next()

		if more {
		}

		return frame.Function, fmt.Sprintf("%v:%v", frame.File, frame.Line)
	}
}

func callers(depth int) []uintptr {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(2+depth, pcs[:])
	var st = pcs[0:n]
	return st
}
