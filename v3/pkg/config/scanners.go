package config

import (
	"encoding"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

var textUnmarshallType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func typeScanner(field *Field) {
	ptrType := field.value.Type()
	fieldType := field.field.Type

	if ptrType.Implements(textUnmarshallType) || fieldType.Implements(textUnmarshallType) {
		field.hasFields = false
	} else {
		fieldKind := fieldType.Kind()

		if fieldKind == reflect.Ptr {
			fieldKind = fieldType.Elem().Kind()
		}
		field.isIgnored = fieldKind == reflect.Interface || fieldKind == reflect.Func
	}
}

type WithDefaults interface {
	SetDefaults()
}

var withDefaultsType = reflect.TypeOf(new(WithDefaults)).Elem()

func defaultsScanner(field *Field) {
	if field.value.Type().Implements(withDefaultsType) {
		field.settings = append(field.settings, setDefaults)
	}
}

func flagTagScanner(field *Field) {
	tag, ok := field.field.Tag.Lookup("flag")

	if !ok {
		return
	}

	var name, description, short string

	parts := strings.Split(tag, ",")

	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		name = parts[0]
		description = parts[1]
	case 3:
		name = parts[0]
		short = parts[1]
		description = parts[2]
	default:
		if len(parts) > 3 {
			name = parts[0]
			short = parts[1]
			description = parts[2]
		}
	}

	if name == "" {
		name = strcase.ToKebab(field.name)
	}

	field.settings = append(field.settings, setFlag(name, short, description))
}

func mapStructureTagScanner(item *Field) {
	tag, ok := item.field.Tag.Lookup("env")

	if !ok {
		return
	}

	if index := strings.Index(tag, ","); index != -1 {
		if tag[:index] == "-" {
			item.isIgnored = true
			return
		}

		squash := strings.Contains(tag[index+1:], "squash")
		if squash && (item.field.Type.Kind() == reflect.Struct || (item.field.Type.Kind() == reflect.Ptr && item.field.Type.Elem().Kind() == reflect.Struct)) {
			item.isSquashed = true
		}

		if tag[:index] != "" {
			item.name = tag[:index]
		}

		return
	}

	if tag == "-" {
		item.isIgnored = true
	} else if tag != "" {
		item.name = tag
	}
}
