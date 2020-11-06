package server

import (
	"emperror.dev/emperror"

	"github.com/vseinstrumentiru/lego/v2/inject"
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

func (c *bootContainer) register(constructor inject.Constructor, options ...inject.RegisterOption) *bootContainer {
	emperror.Panic(c.i.Register(constructor, options...))

	return c
}

func (c *bootContainer) execute(invocation inject.Invocation) *bootContainer {
	emperror.Panic(c.i.Execute(invocation))

	return c
}

func (c *bootContainer) instance(instance interface{}, options ...inject.RegisterOption) *bootContainer {
	emperror.Panic(c.i.Instance(instance, options...))

	return c
}

func (c *bootContainer) make(i inject.Interface) *bootContainer {
	emperror.Panic(c.i.Make(i))

	return c
}
