package env

import (
	"encoding"
	"reflect"
)

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
