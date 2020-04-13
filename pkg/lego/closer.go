package lego

import (
	"emperror.dev/errors"
	"io"
)

func Close(i io.Closer) error {
	if i == nil {
		return nil
	}

	return i.Close()
}

type CloseFn func() error

func (c CloseFn) Close() error {
	return c()
}

type CloserGroup struct {
	closers []io.Closer
}

func (c CloserGroup) Add(closer io.Closer) {
	c.closers = append(c.closers, closer)
}

func (c CloserGroup) Close() error {
	var err error
	for _, fn := range c.closers {
		err = errors.Append(err, fn.Close())
	}

	return err
}
