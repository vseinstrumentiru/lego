package errorhandler

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	logurhandler "emperror.dev/handler/logur"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/errorhandler/sentryprovider"
	"logur.dev/logur"
)

func Provide(p lego2.Process, config Config, logger logur.LoggerFacade) emperror.ErrorHandlerFacade {
	handlers := emperror.ErrorHandlers{}

	if len(config.Providers) == 0 {
		config.Providers = append(config.Providers, LogProvider)
	}

	for _, t := range config.Providers {
		switch t {
		case LogProvider:
			handlers = append(handlers, logurhandler.New(logger))
		case SentryProvider:
			provider, err := sentryprovider.New(p, config.Sentry.DSN)
			emperror.Panic(errors.Wrap(err, "can't load sentry error handler"))
			handlers = append(handlers, provider)
		}
	}

	return handlers
}
