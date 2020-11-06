package container

import (
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/inject"
)

func WithName(name string) inject.RegisterOption {
	return dig.Name(name)
}
