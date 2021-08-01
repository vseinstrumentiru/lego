package handlers

import (
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/log"
)

type Out struct {
	dig.Out
	Handler log.EntryHandler `group:"log.handlers"`
}

func Provider(handler log.EntryHandler) Out {
	return Out{Handler: handler}
}

func Empty() Out {
	return Provider(nil)
}
