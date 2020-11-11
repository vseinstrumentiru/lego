package env

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/vseinstrumentiru/lego/v2/config"
)

func newParser() *parser {
	return &parser{
		configs:          make(map[string]interface{}),
		configTypesCount: make(map[reflect.Type]int),
		validates:        make(map[string]config.Validatable),
	}
}

func parse(v reflect.Value) *parser {
	result := newParser()

	result.scan(v, "", false)

	return result
}

type defaults struct {
	key string
	val config.WithDefaults
}

type parser struct {
	configs          map[string]interface{}
	configTypesCount map[reflect.Type]int
	defaults         []defaults
	validates        map[string]config.Validatable
	keys             []string
}

func (p *parser) isDouble(i interface{}) bool {
	v, ok := normalizeStruct(i)

	if !ok {
		return false
	}

	return p.configTypesCount[v.Type()] > 1
}

func (p *parser) exist(i interface{}) bool {
	v, ok := normalizeStruct(i)

	if !ok {
		return false
	}

	return p.configTypesCount[v.Type()] > 0
}

func (p *parser) scan(v reflect.Value, key string, isSquash bool) {
	{
		var ok bool
		v, ok = normalizeStructValue(v)

		if !ok {
			if key != "" {
				p.keys = append(p.keys, key)
			}
			return
		}
	}

	typ := v.Type()
	ptrI := v.Addr().Interface()
	isAnonymousStructure := typ.Name() == "" || isSquash

	if !isAnonymousStructure {
		configKey := key
		if key == "" {
			configKey = "cfg"
		}
		p.configs[configKey] = ptrI
		p.configTypesCount[typ]++
	}

	if val, ok := ptrI.(config.WithDefaults); ok {
		p.defaults = append(p.defaults, defaults{key: key, val: val})
	}

	if val, ok := ptrI.(config.Validatable); ok {
		p.validates[key] = val
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		t := v.Type().Field(i)
		name, ok := getFieldName(f, t)

		if !ok {
			continue
		}

		nextKey := key
		if name != "" {
			if nextKey != "" {
				nextKey += "."
			}
			nextKey += strcase.ToLowerCamel(name)
		}

		val, ok := normalizeStructValue(f)

		if !ok && f.Kind() == reflect.Ptr && f.IsNil() {
			val = reflect.New(f.Type().Elem())
			f.Set(val)
		}

		p.scan(val, nextKey, name == "")
	}
}

func getFieldName(field reflect.Value, typ reflect.StructField) (string, bool) {
	if !field.CanInterface() {
		return "", false
	}

	if typ.Anonymous && !field.CanSet() {
		return "", false
	}

	loadTag := typ.Tag.Get("load")
	if loadTag == "true" {
		return typ.Name, true
	}

	mapTag := typ.Tag.Get("mapstructure")
	if index := strings.Index(mapTag, ","); index != -1 {
		if mapTag[:index] == "-" {
			return "", false
		}

		squash := strings.Contains(mapTag[index+1:], "squash")
		if squash && typ.Type.Kind() == reflect.Struct {
			return "", true
		}

		return mapTag[:index], true
	} else if len(mapTag) > 0 {
		if mapTag == "-" {
			return "", false
		}
		return mapTag, true
	}

	return typ.Name, true
}

func normalizeStruct(i interface{}) (reflect.Value, bool) {
	return normalizeStructValue(reflect.ValueOf(i))
}

func normalizeStructValue(v reflect.Value) (reflect.Value, bool) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() == reflect.Ptr {
		if !v.IsValid() {
			return reflect.Value{}, false
		}

		if v.IsNil() {
			return reflect.Value{}, false
		}

		if !v.CanInterface() && v.Type().Name() != "" {
			return reflect.Value{}, false
		}

		v = v.Elem()

		if !v.CanAddr() {
			return reflect.Value{}, false
		}
	}

	if !v.IsValid() || v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	return v, true
}
