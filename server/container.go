package server

import (
	"emperror.dev/emperror"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/internal/container"
)

func newContainer() *bootContainer {
	return &bootContainer{
		i: container.New(),
	}
}

type bootContainer struct {
	i container.Container
}

func (c *bootContainer) register(constructor interface{}, options ...dig.ProvideOption) *bootContainer {
	emperror.Panic(c.i.Register(constructor, options...))

	return c
}

func (c *bootContainer) execute(function interface{}) *bootContainer {
	emperror.Panic(c.i.Execute(function))

	return c
}

func (c *bootContainer) instance(instance interface{}, options ...dig.ProvideOption) *bootContainer {
	emperror.Panic(c.i.Instance(instance, options...))

	return c
}

func (c *bootContainer) make(i interface{}) *bootContainer {
	emperror.Panic(c.i.Make(i))

	return c
}
