package sentryprovider

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"time"
)

// Handler is responsible for sending errors to Sentry.
type SentryHandler struct{}

// New creates a new handler.
func New(dsn string) (*SentryHandler, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: false,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create raven client")
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
