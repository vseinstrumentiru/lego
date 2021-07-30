package container

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func New() Container {
	instance := dig.New()

	c := &container{
		di: instance,
	}

	return c
}

type container struct {
	di *dig.Container
}

func (c *container) Visualize(w io.Writer, opts ...dig.VisualizeOption) error {
	return dig.Visualize(c.di, w, opts...)
}

func (c *container) Register(constructor interface{}, options ...dig.ProvideOption) error {
	return errors.WithStack(c.di.Provide(constructor, options...))
}

func (c *container) Instance(instance interface{}, options ...dig.ProvideOption) error {
	t := reflect.ValueOf(instance)

	if err := checkStruct(t); err != nil {
		return errors.WithStack(err)
	}

	funcType := reflect.FuncOf(nil, []reflect.Type{t.Type()}, false)
	f := reflect.MakeFunc(funcType, instanceFn(t))

	return c.Register(f.Interface(), options...)
}

func (c *container) Execute(function interface{}) error {
	return errors.WithStack(c.di.Invoke(function))
}

func (c *container) Make(i interface{}) error {
	val := reflect.ValueOf(i)

	if err := checkStruct(val); err != nil {
		return err
	}
	// setup providers from struct
	if appWithProviders, ok := i.(interface{ Providers() []interface{} }); ok {
		constructors := appWithProviders.Providers()
		for i := 0; i < len(constructors); i++ {
			if err := c.Register(constructors[i]); err != nil {
				return err
			}
		}
	}

	// prepare to make changes to original struct
	ptr := val
	if val.Kind() != reflect.Ptr {
		ptr = reflect.New(val.Type())
		ptr.Elem().Set(val)
		val = ptr.Elem()
	}

	// resolve constructors and collect configuration methods
	t := val.Type()
	var configurations []interface{}

	if appWithConfiguration, ok := i.(interface{ Configurations() []interface{} }); ok {
		configurations = append(configurations, appWithConfiguration.Configurations()...)
	}

	if appWithConfiguration, ok := i.(interface{ With() []interface{} }); ok {
		configurations = append(configurations, appWithConfiguration.With()...)
	}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		if m.Name == "Providers" || m.Name == "Configurations" || m.Name == "With" {
			continue
		}

		if m.Name == "Provide" || strings.HasPrefix(m.Name, "Provide") {
			constructor := val.MethodByName(m.Name).Interface()
			if err := c.Register(constructor); err != nil {
				return err
			}
		} else if m.Name == "Configure" || strings.HasPrefix(m.Name, "Configure") || strings.HasPrefix(m.Name, "With") {
			configurations = append(configurations, val.MethodByName(m.Name).Interface())
		}
	}

	// resolve structure
	if err := c.resolve(ptr); err != nil {
		return errors.WithStack(err)
	}

	// execute configuration methods
	for i := 0; i < len(configurations); i++ {
		if err := c.Execute(configurations[i]); err != nil {
			return err
		}
	}

	return nil
}

func (c *container) resolve(val reflect.Value) error {
	var inFields, outFields []reflect.StructField
	t := val.Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous || !val.Elem().Field(i).CanSet() {
			continue
		}

		if alias, ok := field.Tag.Lookup("name"); !ok {
			inFields = append(inFields, reflect.StructField{
				Name: field.Name,
				Type: field.Type,
				Tag:  reflect.StructTag(fmt.Sprintf(`name:"%s"`, alias)),
			})
		} else {
			inFields = append(inFields, reflect.StructField{
				Name: field.Name,
				Type: field.Type,
				Tag:  `optional:"true"`,
			})
		}

		outFields = append(outFields, reflect.StructField{
			Name: field.Name,
			Type: field.Type,
		})
	}

	if len(inFields) == 0 {
		return nil
	}

	inFields = append(inFields, reflect.StructField{
		Name:      "In",
		Type:      reflect.TypeOf(dig.In{}),
		Anonymous: true,
	})

	outFields = append(outFields, reflect.StructField{
		Name:      "Out",
		Type:      reflect.TypeOf(dig.Out{}),
		Anonymous: true,
	})

	in := reflect.Indirect(reflect.New(reflect.StructOf(inFields)))
	out := reflect.Indirect(reflect.New(reflect.StructOf(outFields)))
	//nolint:lll
	fn := reflect.MakeFunc(reflect.FuncOf([]reflect.Type{in.Type()}, []reflect.Type{out.Type()}, false), func(args []reflect.Value) (results []reflect.Value) {
		arg := args[0]

		for i := 0; i < arg.Type().NumField(); i++ {
			field := arg.Type().Field(i)
			if field.Anonymous {
				continue
			}
			out.FieldByName(field.Name).Set(arg.FieldByName(field.Name))
		}

		return []reflect.Value{out}
	})

	if err := c.di.Invoke(fn.Interface()); err != nil {
		return err
	}

	for i := 0; i < out.Type().NumField(); i++ {
		field := out.Type().Field(i)
		if field.Anonymous {
			continue
		}
		val.Elem().FieldByName(field.Name).Set(out.FieldByName(field.Name))
	}

	return nil
}

func instanceFn(i reflect.Value) func(args []reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{i}
	}
}

func checkStruct(t reflect.Value) error {
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
