package sentryprovider

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/vseinstrumentiru/lego/pkg/lego"
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

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            p.IsDebug(),
		AttachStacktrace: true,
		ServerName:       serverName,
		Environment:      p.Env(),
		Release:          p.Build().Version + "@" + p.Build().CommitHash,
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
