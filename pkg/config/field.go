package config

import (
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type fieldSetting func(field *Field, v *viper.Viper, flagSet *pflag.FlagSet) error

// Field is a part of any structure
type Field struct {
	path       string
	name       string
	field      reflect.StructField
	value      reflect.Value
	isSquashed bool
	isIgnored  bool
	hasFields  bool
	settings   []fieldSetting
}

func newField(path string, field reflect.StructField, value reflect.Value) (*Field, error) {
	value, err := normalizeStructure(value)

	if err != nil {
		return nil, err
	}

	f := &Field{
		path:      path,
		name:      strcase.ToLowerCamel(field.Name),
		value:     value,
		field:     field,
		hasFields: value.Elem().Type().Kind() == reflect.Struct && value.Elem().Type().NumField() > 0,
	}

	return f, nil
}

func (f *Field) key() string {
	if f.path == "" {
		return f.name
	}

	return f.path + "." + f.name
}

func (f *Field) setUp(v *viper.Viper, flagSet *pflag.FlagSet) error {
	for _, fn := range f.settings {
		if err := fn(f, v, flagSet); err != nil {
			return err
		}
	}

	return nil
}
