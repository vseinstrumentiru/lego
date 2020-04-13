package errorhandler

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	logurhandler "emperror.dev/handler/logur"
	"github.com/vseinstrumentiru/lego/internal/monitor/errorhandler/sentryprovider"
	"logur.dev/logur"
)

func Provide(config Config, logger logur.LoggerFacade) emperror.ErrorHandlerFacade {
	handlers := emperror.ErrorHandlers{}

	if len(config.Providers) == 0 {
		config.Providers = append(config.Providers, LogProvider)
	}

	for _, t := range config.Providers {
		switch t {
		case LogProvider:
			handlers = append(handlers, logurhandler.New(logger))
		case SentryProvider:
			provider, err := sentryprovider.New(config.Sentry.DSN)
			emperror.Panic(errors.Wrap(err, "can't load sentry error handler"))
			handlers = append(handlers, provider)
		}
	}

	return handlers
}
