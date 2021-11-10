package config

import (
	"reflect"

	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type configureCallback func(field *Field) error

func noopConfigureCallback(_ *Field) error {
	return nil
}

func configure(cfg interface{}, viper *viper.Viper, flagSet *pflag.FlagSet, cb configureCallback) error {
	v := reflect.ValueOf(cfg)

	if v.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	v, err := normalizeStructure(v)

	if err != nil {
		return err
	}

	if v.Elem().Kind() != reflect.Struct {
		return ErrNotStructure
	}

	scanner := getDefaultScanner()

	fields, err := scan("", v, scanner)

	if err != nil {
		return err
	}

	index := make(map[string]*Field)

	if cb == nil {
		cb = noopConfigureCallback
	}

	for j := 0; j < len(fields); j++ {
		field := fields[j]

		if err = field.setUp(viper, flagSet); err != nil {
			return err
		}

		if _, ok := index[field.key()]; ok {
			return errors.WithDetails(ErrDuplicateKey, "key", field.key())
		}
		index[field.key()] = field

		if err = cb(field); err != nil {
			return err
		}
	}

	if i, ok := cfg.(WithDefaults); ok {
		i.SetDefaults()
	}

	return err
}

func normalizeStructure(elem reflect.Value) (reflect.Value, error) {
	if !elem.IsValid() {
		return elem, ErrNotValidValue
	}

	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			if !elem.CanSet() {
				return elem, ErrNotValidValue
			}

			v := reflect.New(elem.Type().Elem())
			elem.Set(v)
		}
		elem = elem.Elem()

		if !elem.IsValid() {
			return elem, ErrNotValidValue
		}
	}

	if !elem.CanAddr() {
		return elem, ErrNotValidValue
	}

	ptr := elem.Addr()

	if !elem.CanSet() {
		return ptr, ErrNotValidValue
	}

	return ptr, nil
}
