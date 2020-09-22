package exporters

import (
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
)

type Opencensus struct {
	Addr            string
	Insecure        bool
	ReconnectPeriod time.Duration
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
