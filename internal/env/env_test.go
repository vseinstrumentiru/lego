package env

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"

	"github.com/vseinstrumentiru/lego/config"
)

type SubStruct struct{ S1 string }

type structWithDefaults struct{}

func (structWithDefaults) SetDefaults(e config.Env) {}

type structWithValidates struct{}

func (structWithValidates) Validate() error { return nil }

func preloadEnvs() {
	_ = os.Setenv("TEST_A", "test-a")
	_ = os.Setenv("TEST_B", "1")
	_ = os.Setenv("TEST_C_A", "test-c-a")
	_ = os.Setenv("TEST_D_D_A", "2")
	_ = os.Setenv("TEST_D_D_B_A", "test-d-b-a")
	_ = os.Setenv("TEST_D_D_B_B", "3")
	_ = os.Setenv("TEST_F_A", "test-f-a")
	_ = os.Setenv("TEST_F_B", "4")
}

type TestRootConfig struct {
	A string
	B int
	C TestSub1Config
	D struct {
		D TestSub2Config
	}
	E string
	F *TestSub1Config
}

type TestSub1Config struct {
	A string
	B int
}

type TestSub2Config struct {
	A int
	B TestSubSubConfig
	C string
}

type TestSubSubConfig struct {
	A string
	B int
	C string
}

func TestConfigure(t *testing.T) {
	preloadEnvs()
	type args struct {
		Config Config
	}
	tests := []struct {
		name    string
		env     *configEnv
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			"success",
			NewConfigEnv("test").(*configEnv),
			args{&TestRootConfig{}},
			false,
			&TestRootConfig{
				A: "test-a",
				B: 1,
				C: TestSub1Config{
					A: "test-c-a",
				},
				D: struct{ D TestSub2Config }{D: TestSub2Config{
					A: 2,
					B: TestSubSubConfig{
						A: "test-d-b-a",
						B: 3,
					},
				}},
				F: &TestSub1Config{
					A: "test-f-a",
					B: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.env.Configure(tt.args.Config); (err != nil) != tt.wantErr {
				t.Errorf("Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.Config, tt.want) {
				t.Errorf("normalizeStruct() got = %v, want %v", tt.args.Config, tt.want)
			}
		})
	}
}

func Test_normalizeStruct(t *testing.T) {
	type testStruct struct {
		Val int
	}

	val := testStruct{Val: 1}

	tests := []struct {
		name  string
		args  interface{}
		want  interface{}
		want1 bool
	}{
		{"success_no_pointer", testStruct{1}, val, true},
		{"success_pointer", &testStruct{1}, val, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := normalizeStruct(tt.args)
			if !reflect.DeepEqual(got.Interface(), tt.want) {
				t.Errorf("normalizeStruct() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("normalizeStruct() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_normalizeStructValue(t *testing.T) {
	type testStruct struct {
		Val int
	}

	val := testStruct{Val: 1}
	primitivePtr := &val.Val

	tests := []struct {
		name  string
		args  reflect.Value
		want  interface{}
		want1 bool
	}{
		{"success_no_pointer", reflect.ValueOf(testStruct{1}), val, true},
		{"success_pointer", reflect.ValueOf(&testStruct{1}), val, true},
		{"not_ok_primitive_1", reflect.ValueOf(5), reflect.Value{}, false},
		{"not_ok_primitive_2", reflect.ValueOf(primitivePtr), reflect.Value{}, false},
		{"not_ok_nil_1", reflect.ValueOf((*testStruct)(nil)), reflect.Value{}, false},
		{"not_ok_nil_2", reflect.ValueOf(nil), reflect.Value{}, false},
		{"not_ok_nil_2", reflect.Indirect(reflect.ValueOf((*testStruct)(nil)).Elem()), reflect.Value{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := normalizeStructValue(tt.args)

			if got1 != tt.want1 {
				t.Errorf("normalizeStructValue() got1 = %v, want %v", got1, tt.want1)
			}

			if got1 {
				assert.Equal(t, tt.want, got.Interface(), "normalizeStructValue() got = %v, want %v", got.Interface(), tt.want)
			} else {
				assert.Equal(t, tt.want, got, "normalizeStructValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want *parser
	}{
		{"defaults",
			&struct{ A structWithDefaults }{},
			&parser{
				configs: map[string]interface{}{
					"a": &structWithDefaults{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(structWithDefaults{}): 1,
				},
				defaults: []defaults{
					{"a", &structWithDefaults{}},
				},
				validates: map[string]config.Validatable{},
				keys:      nil,
			},
		},
		{"defaults_pointer",
			&struct{ A *structWithDefaults }{},
			&parser{
				configs: map[string]interface{}{
					"a": &structWithDefaults{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(structWithDefaults{}): 1,
				},
				defaults: []defaults{
					{"a", &structWithDefaults{}},
				},
				validates: map[string]config.Validatable{},
				keys:      nil,
			},
		},
		{"validates",
			&struct{ A structWithValidates }{},
			&parser{
				configs: map[string]interface{}{
					"a": &structWithValidates{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(structWithValidates{}): 1,
				},
				defaults: nil,
				validates: map[string]config.Validatable{
					"a": &structWithValidates{},
				},
				keys: nil,
			},
		},
		{"validates_pointer",
			&struct{ A *structWithValidates }{},
			&parser{
				configs: map[string]interface{}{
					"a": &structWithValidates{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(structWithValidates{}): 1,
				},
				defaults: nil,
				validates: map[string]config.Validatable{
					"a": &structWithValidates{},
				},
				keys: nil,
			},
		},
		{"unexported",
			&struct {
				a1  structWithDefaults
				a2  *structWithDefaults
				a3  structWithValidates
				a4  *structWithValidates
				a5  string
				a6  *string
				A7  string     `mapstructure:"-"`
				A8  *string    `mapstructure:"-"`
				A9  SubStruct  `mapstructure:"-"`
				A10 *SubStruct `mapstructure:"-"`
				A11 struct{}
				_   SubStruct
				_   *SubStruct
			}{},
			&parser{
				configs:          map[string]interface{}{},
				configTypesCount: map[reflect.Type]int{},
				defaults:         nil,
				validates:        map[string]config.Validatable{},
				keys:             nil,
			},
		},
		{"fields",
			&struct {
				S1 string
				S2 *string
				A1 struct {
					B1 string
				}
				A2 *struct {
					B1 string
				}
				A3 SubStruct
				A4 *SubStruct
				A5 struct {
					SubStruct
					S1 string
				}
				A6 struct {
					*SubStruct
					S1 string
				}
			}{},
			&parser{
				configs: map[string]interface{}{
					"a3":           &SubStruct{},
					"a4":           &SubStruct{},
					"a5.subStruct": &SubStruct{},
					"a6.subStruct": &SubStruct{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(SubStruct{}): 4,
				},
				defaults:  nil,
				validates: map[string]config.Validatable{},
				keys: []string{
					"s1",
					"s2",
					"a1.b1",
					"a2.b1",
					"a3.s1",
					"a4.s1",
					"a5.subStruct.s1",
					"a5.s1",
					"a6.subStruct.s1",
					"a6.s1",
				},
			},
		},
		{"embedded_fields",
			&struct {
				A1 struct {
					SubStruct
					structWithDefaults
					structWithValidates
				}
				A2 struct {
					*SubStruct
					*structWithDefaults
					*structWithValidates
				}
				A3 struct {
					SubStruct `mapstructure:"b1"`
				}
				A4 struct {
					*SubStruct `mapstructure:"b1"`
				}
				A5 struct {
					SubStruct `mapstructure:",squash"`
				}
				A6 struct {
					*SubStruct `mapstructure:",squash"`
				}
			}{},
			&parser{
				configs: map[string]interface{}{
					"a1.subStruct": &SubStruct{},
					"a2.subStruct": &SubStruct{},
					"a3.b1":        &SubStruct{},
					"a4.b1":        &SubStruct{},
				},
				configTypesCount: map[reflect.Type]int{
					reflect.TypeOf(SubStruct{}): 4,
				},
				defaults: []defaults{
					{"a1", &struct {
						SubStruct
						structWithDefaults
						structWithValidates
					}{}},
					{"a2", &struct {
						*SubStruct
						*structWithDefaults
						*structWithValidates
					}{SubStruct: &SubStruct{}}},
				},
				validates: map[string]config.Validatable{
					"a1": &struct {
						SubStruct
						structWithDefaults
						structWithValidates
					}{},
					"a2": &struct {
						*SubStruct
						*structWithDefaults
						*structWithValidates
					}{SubStruct: &SubStruct{}},
				},
				keys: []string{
					"a1.subStruct.s1",
					"a2.subStruct.s1",
					"a3.b1.s1",
					"a4.b1.s1",
					"a5.s1",
					"a6.s1",
				},
			},
		},
	}
	opts := []cmp.Option{
		cmp.AllowUnexported(sync.Mutex{}, defaults{}),
		cmp.AllowUnexported(struct {
			SubStruct
			structWithDefaults
			structWithValidates
		}{}, struct {
			*SubStruct
			*structWithDefaults
			*structWithValidates
		}{}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(reflect.ValueOf(tt.args))
			assert.DeepEqual(t, tt.want.configs, got.configs, opts...)
			assert.DeepEqual(t, tt.want.configTypesCount, got.configTypesCount, opts...)
			assert.DeepEqual(t, tt.want.defaults, got.defaults, opts...)
			assert.DeepEqual(t, tt.want.validates, got.validates, opts...)
			assert.DeepEqual(t, tt.want.keys, got.keys, opts...)
		})
	}
}

func Test_parsed_exist(t *testing.T) {
	type fields struct {
		configs          map[string]interface{}
		configTypesCount map[reflect.Type]int
		defaults         []defaults
		validates        map[string]config.Validatable
		keys             []string
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				configs:          tt.fields.configs,
				configTypesCount: tt.fields.configTypesCount,
				defaults:         tt.fields.defaults,
				validates:        tt.fields.validates,
				keys:             tt.fields.keys,
			}
			if got := p.exist(tt.args.i); got != tt.want {
				t.Errorf("exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsed_isDouble(t *testing.T) {
	type fields struct {
		configs          map[string]interface{}
		configTypesCount map[reflect.Type]int
		defaults         []defaults
		validates        map[string]config.Validatable
		keys             []string
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				configs:          tt.fields.configs,
				configTypesCount: tt.fields.configTypesCount,
				defaults:         tt.fields.defaults,
				validates:        tt.fields.validates,
				keys:             tt.fields.keys,
			}
			if got := p.isDouble(tt.args.i); got != tt.want {
				t.Errorf("isDouble() = %v, want %v", got, tt.want)
			}
		})
	}
}
