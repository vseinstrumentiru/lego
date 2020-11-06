package container

import di "github.com/vseinstrumentiru/lego/v2/inject"

type Container interface {
	Register(constructor di.Constructor, options ...di.RegisterOption) error
	Instance(instance interface{}, options ...di.RegisterOption) error
	Execute(invocation di.Invocation) error
	Make(i di.Interface) error
}
