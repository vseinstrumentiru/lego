package config

import (
	"reflect"
	"strings"
)

func parseConfig(config interface{}, baseKey string) (defaults map[string]ConfigWithDefaults, validates map[string]Validateable) {
	defaults = make(map[string]ConfigWithDefaults)
	validates = make(map[string]Validateable)
	configParser(config, baseKey, defaults, validates)

	return
}

func configParser(config interface{}, key string, defaults map[string]ConfigWithDefaults, validates map[string]Validateable) {
	if key == "" {
		return
	}

	if val, ok := config.(ConfigWithDefaults); ok {
		defaults[key] = val
	}

	if val, ok := config.(Validateable); ok {
		validates[key] = val
	}

	v := reflect.ValueOf(config)

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name, ok := getFieldName(v.Type().Field(i))

		if !ok || !f.CanInterface() {
			continue
		}

		nextKey := key
		if name != "" {
			nextKey += "." + name
		}

		val := f.Interface()

		if f.Kind() == reflect.Ptr {
			ptr := f
			if f.IsNil() {
				ptr = reflect.New(f.Type().Elem())
			}

			val = reflect.Indirect(ptr).Interface()
		}

		configParser(val, nextKey, defaults, validates)
	}
}

func getFieldName(f reflect.StructField) (string, bool) {
	tagValue := f.Tag.Get("mapstructure")

	if index := strings.Index(tagValue, ","); index != -1 {
		if tagValue[:index] == "-" {
			return "", false
		}

		// If "squash" is specified in the tag, we squash the field down.
		squash := strings.Index(tagValue[index+1:], "squash") != -1
		if squash && f.Type.Kind() == reflect.Struct {
			return "", true
		}

		return tagValue[:index], true
	} else if len(tagValue) > 0 {
		if tagValue == "-" {
			return "", false
		}
		return tagValue, true
	}

	return f.Name, true
}
