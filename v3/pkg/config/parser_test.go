package config_test

import (
	"os"
	"testing"

	"emperror.dev/errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/vseinstrumentiru/lego/v3/pkg/config"
)

func TestParser_Unmarshal_Flags(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string `flag:""`
	}
	v := &root{}
	p := config.New(config.WithArgs([]string{"", "--field-1=arg", "--wrong=wrong"}))
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "arg", v.Field1)
}

func TestParser_Unmarshal_FlagsWithEnvs(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string `flag:""`
		Field2 string
	}
	v := &root{}
	assert.NoError(t, os.Setenv("FIELD1", "fail"))
	assert.NoError(t, os.Setenv("FIELD2", "success"))
	p := config.New(config.WithArgs([]string{"", "--field-1=arg", "--wrong=wrong"}))
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "arg", v.Field1)
	assert.Equal(t, "success", v.Field2)
}

func TestParser_Unmarshal_EnvsWithPrefix(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string
		Field2 string
	}
	v := &root{}
	assert.NoError(t, os.Setenv("TEST_FIELD1", "success"))
	assert.NoError(t, os.Setenv("FIELD2", "fail"))
	p := config.New(config.WithEnvPrefix("test"))
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "success", v.Field1)
	assert.Equal(t, "", v.Field2)
}

func TestParser_Unmarshal_SubstructureEnvs(t *testing.T) {
	os.Clearenv()
	type Substructure struct {
		SubField1 string
	}

	type root struct {
		Field1       string
		Field2       Substructure
		Substructure `env:"sub"`
		Field3       Substructure `env:",squash"`
	}
	v := &root{}
	assert.NoError(t, os.Setenv("FIELD1", "success_1"))
	assert.NoError(t, os.Setenv("FIELD2_SUBFIELD1", "success_2_1"))
	assert.NoError(t, os.Setenv("SUB_SUBFIELD1", "success_sub_1"))
	assert.NoError(t, os.Setenv("SUBSTRUCTURE_SUBFIELD1", "fail_sub_1"))
	assert.NoError(t, os.Setenv("SUBFIELD1", "success_subfield1"))
	p := config.New()
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "success_1", v.Field1)
	assert.Equal(t, "success_2_1", v.Field2.SubField1)
	assert.Equal(t, "success_sub_1", v.Substructure.SubField1)
	assert.Equal(t, "success_subfield1", v.Field3.SubField1)
}

type structWithDefaults struct {
	SubField1 string
}

func (s *structWithDefaults) SetDefaults() {
	s.SubField1 = "default"
}

func TestParser_Unmarshal_Defaults(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string
		Field2 structWithDefaults
	}

	v := &root{}
	assert.NoError(t, os.Setenv("FIELD1", "success_1"))
	p := config.New()
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "success_1", v.Field1)
	assert.Equal(t, "default", v.Field2.SubField1)
}

func TestParser_Unmarshal_DefaultsOverride(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string
		Field2 structWithDefaults
	}

	v := &root{}
	assert.NoError(t, os.Setenv("FIELD1", "success_1"))
	assert.NoError(t, os.Setenv("FIELD2_SUBFIELD1", "success_2_1"))
	p := config.New()
	err := p.Unmarshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "success_1", v.Field1)
	assert.Equal(t, "success_2_1", v.Field2.SubField1)
}

func TestParser_Unmarshal_ValidatorFails(t *testing.T) {
	os.Clearenv()
	type root struct {
		Field1 string `validate:"required"`
		Field2 int    `validate:"gte=10"`
		Field3 struct {
			Field1 string `validate:"required"`
		}
	}

	v := &root{}
	p := config.New()
	err := p.Unmarshal(v)
	assert.True(t, err != nil)
	errs := errors.GetErrors(err)
	assert.Len(t, errs, 3)
	assert.ErrorAs(t, errs[0], new(validator.FieldError))
	fieldErr := errs[0].(validator.FieldError)
	assert.Equal(t, "root.Field1", fieldErr.Namespace())
	assert.ErrorAs(t, errs[1], new(validator.FieldError))
	fieldErr = errs[1].(validator.FieldError)
	assert.Equal(t, "root.Field2", fieldErr.Namespace())
	assert.ErrorAs(t, errs[2], new(validator.FieldError))
	fieldErr = errs[2].(validator.FieldError)
	assert.Equal(t, "root.Field3.Field1", fieldErr.Namespace())
}
