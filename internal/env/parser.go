package env

import (
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"

	"github.com/vseinstrumentiru/lego/v2/config"
)

func newParser() *parser {
	return &parser{
		configs:   make(map[reflect.Type]*Instances),
		validates: make(map[string]config.Validatable),
	}
}

type defaults struct {
	key string
	val config.WithDefaults
}

type Instance struct {
	Val       interface{}
	IsDefault bool
	Key       string
}

type Instances struct {
	DefaultKey int
	Items      []Instance
}

func (i *Instances) add(v interface{}, key string, isDefault bool) error {
	i.Items = append(i.Items, Instance{
		Val:       v,
		IsDefault: isDefault,
		Key:       key,
	})

	if isDefault {
		if len(i.Items) > 1 && i.Items[i.DefaultKey].IsDefault {
			return errors.New("two default configs with same type")
		}

		i.DefaultKey = len(i.Items) - 1
	}

	return nil
}

type flags struct {
	isSquash  bool
	isDefault bool
}

type parser struct {
	configs   map[reflect.Type]*Instances
	defaults  []defaults
	validates map[string]config.Validatable
	keys      []string
}

func (p *parser) isDouble(i interface{}) bool {
	return false
}

func (p *parser) exist(i interface{}) bool {
	v, ok := normalizeStruct(i)

	if !ok {
		return false
	}

	_, ok = p.configs[v.Type()]

	return ok
}

func (p *parser) scan(v reflect.Value, key string, f flags) error {
	{
		var ok bool
		v, ok = normalizeStructValue(v)

		if !ok {
			if key != "" {
				p.keys = append(p.keys, key)
			}
			return nil
		}
	}

	typ := v.Type()
	ptrI := v.Addr().Interface()
	isAnonymousStructure := typ.Name() == "" || (f.isSquash && !f.isDefault)

	if !isAnonymousStructure {
		configKey := key
		if key == "" {
			configKey = "cfg"
		}

		i, ok := p.configs[typ]
		if !ok {
			i = &Instances{}
			p.configs[typ] = i
		}
		if err := i.add(ptrI, configKey, f.isDefault); err != nil {
			return err
		}
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
		info := getFieldInfo(f, t)

		if info.Ignore {
			continue
		}

		nextKey := key
		if info.Name != "" {
			if nextKey != "" {
				nextKey += "."
			}
			nextKey += strcase.ToLowerCamel(info.Name)
		}

		val, ok := normalizeStructValue(f)

		if !ok && f.Kind() == reflect.Ptr && f.IsNil() {
			val = reflect.New(f.Type().Elem())
			f.Set(val)
		}

		if err := p.scan(val, nextKey, flags{isSquash: info.IsSquashed, isDefault: info.IsDefault}); err != nil {
			return err
		}
	}

	return nil
}

type FieldInfo struct {
	Name       string
	Ignore     bool
	IsDefault  bool
	IsSquashed bool
}

func getFieldInfo(field reflect.Value, typ reflect.StructField) FieldInfo {
	if !field.CanInterface() {
		return FieldInfo{Ignore: true}
	}

	if typ.Anonymous && !field.CanSet() {
		return FieldInfo{Ignore: true}
	}

	f := FieldInfo{Name: typ.Name}

	if defaultTag := typ.Tag.Get("default"); defaultTag == "true" {
		f.IsDefault = true
	}

	loadTag := typ.Tag.Get("load")
	if loadTag == "true" {
		return f
	}

	mapTag := typ.Tag.Get("mapstructure")
	if index := strings.Index(mapTag, ","); index != -1 {
		if mapTag[:index] == "-" {
			return FieldInfo{Ignore: true}
		}

		squash := strings.Contains(mapTag[index+1:], "squash")
		if squash && (typ.Type.Kind() == reflect.Struct || (typ.Type.Kind() == reflect.Ptr && typ.Type.Elem().Kind() == reflect.Struct)) {
			f.Name = ""
			f.IsSquashed = true
			return f
		}

		f.Name = mapTag[:index]
		return f
	} else if len(mapTag) > 0 {
		if mapTag == "-" {
			return FieldInfo{Ignore: true}
		}
		f.Name = mapTag
		return f
	}

	return f
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
