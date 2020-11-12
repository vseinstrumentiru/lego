package env

import (
	"reflect"
	"testing"

	"gotest.tools/assert"
)

func Test_getFieldInfo(t *testing.T) {
	var field struct{}

	type args struct {
		field reflect.Value
		typ   reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want FieldInfo
	}{
		{"simple", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field)}}, FieldInfo{Name: "test"}},
		{"simple/default", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`default:"true"`)}}, FieldInfo{Name: "test", IsDefault: true}},
		{"anonymous", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Anonymous: true}}, FieldInfo{Ignore: true}},
		{"mapstructure/ignore", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-"`)}}, FieldInfo{Ignore: true}},
		{"mapstructure/squash", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:",squash"`)}}, FieldInfo{IsSquashed: true}},
		{"mapstructure/default+squash", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:",squash" default:"true"`)}}, FieldInfo{IsSquashed: true, IsDefault: true}},
		{"mapstructure/name", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"replaced"`)}}, FieldInfo{Name: "replaced"}},
		{"mapstructure/name+squash", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"replaced,squash"`)}}, FieldInfo{IsSquashed: true}},
		{"load", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-" load:"true"`)}}, FieldInfo{Name: "test"}},
		{"load+default", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-" load:"true" default:"true"`)}}, FieldInfo{Name: "test", IsDefault: true}},
		{"ignore+default", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-" default:"true"`)}}, FieldInfo{Ignore: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFieldInfo(tt.args.field, tt.args.typ)
			assert.DeepEqual(t, tt.want, got)
		})
	}
}
