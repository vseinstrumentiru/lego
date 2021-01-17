package container

import (
	"emperror.dev/emperror"
	"go.uber.org/dig"
)

func NewChain(parent Container) ChainContainer {
	c := &chainContainer{
		Parent: parent,
	}

	return c
}

type chainContainer struct {
	Parent Container
}

func (c *chainContainer) Register(constructor interface{}, options ...dig.ProvideOption) ChainContainer {
	emperror.Panic(c.Parent.Register(constructor, options...))

	return c
}

func (c *chainContainer) Execute(function interface{}) ChainContainer {
	emperror.Panic(c.Parent.Execute(function))

	return c
}

func (c *chainContainer) Instance(instance interface{}, options ...dig.ProvideOption) ChainContainer {
	emperror.Panic(c.Parent.Instance(instance, options...))

	return c
}

func (c *chainContainer) Make(i interface{}) ChainContainer {
	emperror.Panic(c.Parent.Make(i))

	return c
}
