package config

import (
	"reflect"
	"testing"
)

func TestApplication_FullName(t *testing.T) {
	type fields struct {
		Name       string
		DataCenter string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"success_1", fields{"name", "dc"}, "name-dc"},
		{"success_2", fields{"name-with_$peci@1s", "dc-with_$peci@1s"}, "name-with_$peci@1s-dc-with_$peci@1s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Application{
				Name:       tt.fields.Name,
				DataCenter: tt.fields.DataCenter,
			}
			if got := c.FullName(); got != tt.want {
				t.Errorf("FullName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUndefined(t *testing.T) {
	tests := []struct {
		name string
		want *Application
	}{
		{"success", &Application{Name: "undefined", DataCenter: "undefined"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UndefinedApplication(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UndefinedApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}
