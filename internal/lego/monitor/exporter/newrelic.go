package exporter

import (
	"emperror.dev/emperror"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vseinstrumentiru/lego/internal/lego"
)

type nrLogErr struct {
	log lego.LogErr
}

func (l nrLogErr) Error(msg string, context map[string]interface{}) {
	l.log.Error(msg, context)
}

func (l nrLogErr) Warn(msg string, context map[string]interface{}) {
	l.log.Warn(msg, context)
}

func (l nrLogErr) Info(msg string, context map[string]interface{}) {
	l.log.Info(msg, context)
}

func (l nrLogErr) Debug(msg string, context map[string]interface{}) {
	l.log.Debug(msg, context)
}

func (nrLogErr) DebugEnabled() bool {
	return false
}

// NewRelicMiddleware it first reads config from env vars, then sets AppName and License key
func NewRelicMiddleware(appName, key string, logErr lego.LogErr) mux.MiddlewareFunc {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(key),
		newrelic.ConfigLogger(nrLogErr{logErr}),
	)

	if err != nil {
		emperror.Panic(err)
	}

	return nrgorilla.Middleware(app)
}
