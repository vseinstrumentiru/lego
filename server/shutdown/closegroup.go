package shutdown

import (
	"io"

	"emperror.dev/errors"
)

var _ io.Closer = &CloseGroup{}

type CloseGroup struct {
	closers []io.Closer
}

func (g *CloseGroup) Close() (err error) {
	for i := 0; i < len(g.closers); i++ {
		err = errors.Append(err, g.closers[i].Close())
	}

	return err
}

func (g *CloseGroup) Add(i io.Closer) {
	g.closers = append([]io.Closer{i}, g.closers...)
}
