package lego

import (
	"github.com/vseinstrumentiru/lego/internal/lego"
	"io"
)

func Close(i io.Closer) error {
	return lego.Close(i)
}

type CloseFn = lego.CloseFn
type CloserGroup = lego.CloserGroup
