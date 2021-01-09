package container

import (
	"go.uber.org/dig"
)

func WithName(name string) dig.ProvideOption {
	return dig.Name(name)
}
