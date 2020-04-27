package sentryprovider

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"os"
	"time"
)

// Handler is responsible for sending errors to Sentry.
type SentryHandler struct{}

// New creates a new handler.
func New(p lego.Process, dsn string) (*SentryHandler, error) {
	serverName := p.Name()
	if p.DataCenterName() != "" {
		serverName += ":" + p.DataCenterName()
	}

	release := p.Build().Version
	if release != "" {
		release += "@"
	}
	release += p.Build().CommitHash

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            p.IsDebug(),
		AttachStacktrace: true,
		ServerName:       p.Name(),
		Environment:      p.Env(),
		Release:          release,
	})

	host, _ := os.Hostname()
	dataCenter := p.DataCenterName()
	branch := p.Build().Version

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("data-center", dataCenter)
		scope.SetTag("branch", branch)
		scope.SetTag("host", host)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create sentry client")
	}

	return &SentryHandler{}, nil
}

func (h *SentryHandler) Handle(err error) {
	if err == error(nil) || err == nil {
		return
	}

	sentry.CaptureException(err)
}

func (h *SentryHandler) Close() error {
	sentry.Flush(2 * time.Second)

	return nil
}
