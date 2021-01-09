package container

import (
	"go.uber.org/dig"
)

type Container interface {
	Register(constructor interface{}, options ...dig.ProvideOption) error
	Instance(instance interface{}, options ...dig.ProvideOption) error
	Execute(function interface{}) error
	Make(i interface{}) error
}
