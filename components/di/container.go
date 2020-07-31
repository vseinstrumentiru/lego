package di

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/goava/di"
	"reflect"
	"strings"
)

func NewContainer() *container {
	di, err := di.New()
	emperror.Panic(err)
	c := &container{
		di: di,
	}

	return c
}

type Constructor = di.Constructor
type Invocation = di.Invocation
type Interface = di.Interface
type Option = di.Option
type RegisterOption = di.ProvideOption
type RegisterParams = di.ProvideParams
type InvokeOption = di.InvokeOption
type MakeOption = di.ResolveOption
type Inject = di.Inject

func As(i ...Interface) RegisterOption {
	return di.As(i...)
}

type container struct {
	di *di.Container
}

func (c *container) Register(constructor Constructor, options ...RegisterOption) {
	err := c.di.Provide(constructor, options...)
	emperror.Panic(err)
}

func (c *container) Instance(instance Interface, options ...RegisterOption) {
	t := reflect.ValueOf(instance)
	emperror.Panic(isStruct(t))

	funcType := reflect.FuncOf(nil, []reflect.Type{t.Type()}, false)
	f := reflect.MakeFunc(funcType, instanceFn(t))

	c.Register(f.Interface(), options...)
}

func (c *container) Invoke(invocation Invocation, options ...InvokeOption) {
	err := c.di.Invoke(invocation, options...)
	emperror.Panic(err)
}

func (c *container) MakeApp(app interface{}) {
	val := reflect.ValueOf(app)
	emperror.Panic(isStruct(val))

	t := val.Type()

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		if m.Name == "Provide" || strings.HasPrefix(m.Name, "Provide") {
			constructor := val.MethodByName(m.Name).Interface()
			c.Register(constructor)
		}
	}

	c.Make(app)
}

func (c *container) Make(i Interface, options ...MakeOption) {
	err := c.di.Resolve(i, options...)
	emperror.Panic(err)
}

// func registerFn(fn reflect.Value) func([]reflect.Value) []reflect.Value {
// 	return func(in []reflect.Value) []reflect.Value {
// 		return fn.Call(in)
// 	}
// }

func instanceFn(i reflect.Value) func(args []reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{i}
	}
}

func isStruct(t reflect.Value) error {
	if t.Kind() == reflect.Ptr {
		if t.IsNil() || !t.IsValid() {
			return errors.New("nil instance presented")
		}

		ti := reflect.Indirect(t)

		if ti.Kind() != reflect.Struct {
			return errors.New("instance must be struct or non-nil pointer to struct")
		}
	} else if t.Kind() != reflect.Struct {
		return errors.New("instance must be struct or non-nil pointer to struct")
	}

	return nil
}
