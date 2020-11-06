package lego

import (
	"io"

	"emperror.dev/errors"

	"github.com/vseinstrumentiru/lego/v2/internal/deprecated"
)

// Deprecated: will deleted in next versions.
type App interface{}

// Deprecated: will deleted in next versions.
type CloseFn func() error

// Deprecated: will deleted in next versions.
func (c CloseFn) Close() error {
	return c()
}

// Deprecated: will deleted in next versions.
type CloserGroup struct {
	closers []io.Closer
}

// Deprecated: will deleted in next versions.
func (c CloserGroup) Add(closer io.Closer) {
	c.closers = append(c.closers, closer)
}

// Deprecated: will deleted in next versions.
func (c CloserGroup) Close() error {
	var err error
	for i := len(c.closers) - 1; i >= 0; i-- {
		err = errors.Append(err, c.closers[i].Close())
	}

	return err
}

// Deprecated: use multilog.Logger
type LogErr = deprecated.LogErr

// Deprecated: use multilog.Logger
type LogErrF = deprecated.LogErrF

// Deprecated: use multilog.Logger
func NewLogErrF(logErr LogErr) LogErrF {
	return deprecated.NewLogErr(logErr)
}

// Deprecated: use DI
type Process = deprecated.Process

// Deprecated: use LeGo V2
type EventManager = deprecated.EventManager

// Deprecated: use LeGo V2
type Config = deprecated.Config
