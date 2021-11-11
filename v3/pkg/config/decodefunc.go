package config

import (
	"encoding"
	"reflect"
)

// StringToTypeDecoder is a mapstructure decoder for any type
// which implements encoding.TextUnmarshaler interface
func StringToTypeDecoder(
	f reflect.Type,
	t reflect.Type,
	data interface{},
) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	val := reflect.New(t)
	i, ok := val.Interface().(encoding.TextUnmarshaler)

	if !ok {
		return data, nil
	}

	err := i.UnmarshalText([]byte(data.(string)))

	return val.Elem().Interface(), err
}
