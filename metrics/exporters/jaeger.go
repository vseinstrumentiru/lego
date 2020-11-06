package exporters

import "emperror.dev/errors"

type Jaeger struct {
	Addr string
}

func (c Jaeger) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("jaeger: addr is empty"))
	}

	return
}
