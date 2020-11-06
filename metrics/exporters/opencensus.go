package exporters

import (
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"emperror.dev/errors"
)

type Opencensus struct {
	Addr            string
	Insecure        bool
	ReconnectPeriod time.Duration
}

func (c *Opencensus) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("opencensus: addr is empty"))
	}

	return
}

func (c *Opencensus) Options() []ocagent.ExporterOption {
	options := []ocagent.ExporterOption{
		ocagent.WithAddress(c.Addr),
		ocagent.WithReconnectionPeriod(c.ReconnectPeriod),
	}

	if c.Insecure {
		options = append(options, ocagent.WithInsecure())
	}

	return options
}
