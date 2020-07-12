package exporters

import "emperror.dev/errors"

type NewRelic struct {
	Enabled          bool
	TelemetryEnabled bool
	Key              string
}

func (c NewRelic) Validate() (err error) {
	if c.Key == "" {
		err = errors.Append(err, errors.New("newrelic: key is empty"))
	}

	return
}
