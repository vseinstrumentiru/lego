package logger

import (
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/log"
)

type args struct {
	dig.In
	App    *config.Application `optional:"true"`
	Config *Config             `optional:"true"`

	Handlers []log.EntryHandler `group:"log.handlers"`
}

func Provide(in args) log.Logger {
	if in.Config == nil {
		level := logur.Error
		if in.App.DebugMode {
			level = logur.Trace
		}

		in.Config = &Config{Level: level}
	}

	opts := make([]log.Option, 0, len(in.Handlers))

	for _, handler := range in.Handlers {
		if handler == nil {
			continue
		}
		opts = append(opts, log.WithHandler(handler))
	}

	logger := log.New(in.Config.Level, opts...)

	if in.App != nil && in.App.Name != config.Undefined && in.App.Name != "" {
		fields := make(map[string]interface{})
		fields["app"] = in.App.Name

		logger = logger.WithFields(fields)
	}

	return logger
}
