package config

import (
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/internal/config/env"
	di "github.com/vseinstrumentiru/lego/internal/container"
)

type RawConfig interface{}

type argsIn struct {
	dig.In
	Config    RawConfig
	Env       env.Env
	Container di.Container
}

func Configure(in argsIn) error {
	v := reflect.ValueOf(in.Config)
	if v.Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type()))
	}

	parsed := parse(in.Config)

	for key, cfg := range parsed.defaults {
		cfg.SetDefaults(in.Env.Sub(key))
	}

	err := in.Env.Load(in.Config, parsed.keys...)

	for _, i := range parsed.validates {
		err = errors.Append(err, i.Validate())
	}

	if err != nil {
		return err
	}

	for key, i := range parsed.configs {
		var err error
		if parsed.configTypesCount[reflect.TypeOf(i)] > 1 {
			err = in.Container.Instance(i, di.WithName(key))
		} else {
			err = in.Container.Instance(i)
		}

		if err != nil {
			return err
		}
	}

	if parsed.configTypesCount[reflect.TypeOf(&config.Application{})] == 0 {
		err = in.Container.Instance(config.Undefined())
	}

	return err
}

func parse(i interface{}) *parsed {
	result := &parsed{
		configs:          make(map[string]interface{}),
		configTypesCount: make(map[reflect.Type]int),
		defaults:         make(map[string]config.ConfigWithDefaults),
		validates:        make(map[string]config.Validateable),
	}

	scan(i, "", result)

	return result
}

type parsed struct {
	configs          map[string]interface{}
	configTypesCount map[reflect.Type]int
	defaults         map[string]config.ConfigWithDefaults
	validates        map[string]config.Validateable
	keys             []string
}

func scan(i interface{}, key string, p *parsed) {
	v := reflect.ValueOf(i)

	configKey := "cfg." + key
	if key == "" {
		configKey = "cfg"
	}

	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		p.configs[configKey] = i
		p.configTypesCount[v.Type()]++
		v = v.Elem()
	} else if v.Kind() == reflect.Struct {
		p.configs[configKey] = v.Addr().Interface()
		p.configTypesCount[v.Addr().Type()]++
	} else {
		if key != "" {
			p.keys = append(p.keys, key)
		}
		return
	}

	if val, ok := i.(config.ConfigWithDefaults); ok {
		p.defaults[key] = val
	}

	if val, ok := i.(config.Validateable); ok {
		p.validates[key] = val
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name, ok := getFieldName(v.Type().Field(i))

		if !ok || !f.CanInterface() {
			continue
		}

		nextKey := key
		if name != "" {
			if nextKey != "" {
				nextKey += "."
			}
			nextKey += strcase.ToLowerCamel(name)
		}

		var val interface{}

		if f.Kind() == reflect.Ptr {
			ptr := f
			if f.IsNil() {
				ptr = reflect.New(f.Type().Elem())
				f.Set(ptr)
			}
			val = ptr.Interface()
		} else {
			val = f.Addr().Interface()
		}

		scan(val, nextKey, p)
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
