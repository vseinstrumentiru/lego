package container

import (
	"io"

	"go.uber.org/dig"
)

type Container interface {
	Register(constructor interface{}, options ...dig.ProvideOption) error
	Instance(instance interface{}, options ...dig.ProvideOption) error
	Execute(function interface{}) error
	Make(i interface{}) error
	Visualize(w io.Writer, opts ...dig.VisualizeOption) error
}

type ChainContainer interface {
	Register(constructor interface{}, options ...dig.ProvideOption) ChainContainer
	Execute(function interface{}) ChainContainer
	Instance(instance interface{}, options ...dig.ProvideOption) ChainContainer
	Make(i interface{}) ChainContainer
}
