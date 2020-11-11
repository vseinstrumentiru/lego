package env

import (
	"reflect"
	"testing"
)

func Test_getFieldName(t *testing.T) {
	var field struct{}

	type args struct {
		field reflect.Value
		typ   reflect.StructField
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"simple", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field)}}, "test", true},
		{"anonymous", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Anonymous: true}}, "", false},
		{"mapstructure/ignore", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-"`)}}, "", false},
		{"mapstructure/squash", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:",squash"`)}}, "", true},
		{"mapstructure/name", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"replaced"`)}}, "replaced", true},
		{"mapstructure/name+squash", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"replaced,squash"`)}}, "", true},
		{"load/true", args{reflect.ValueOf(field), reflect.StructField{Name: "test", Type: reflect.TypeOf(field), Tag: reflect.StructTag(`mapstructure:"-" load:"true"`)}}, "test", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getFieldName(tt.args.field, tt.args.typ)
			if got != tt.want {
				t.Errorf("getFieldName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getFieldName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
