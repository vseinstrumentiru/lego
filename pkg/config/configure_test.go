package config

import (
	"os"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func itemsToMapCallback(result map[string]*Field) configureCallback {
	return func(item *Field) error {
		result[item.key()] = item

		return nil
	}
}

func TestConfigure_ErrorsNotPointer(t *testing.T) {
	type root struct {
		Field1 string
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	err := configure(root{}, v, flagSet, nil)

	assert.ErrorIs(t, err, ErrNotPointer)
}

func TestConfigure_ErrorsNotStruct(t *testing.T) {
	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	str := "asd"
	err := configure(&str, v, flagSet, nil)

	assert.ErrorIs(t, err, ErrNotStructure)
}

func TestConfigure_ErrorsNotValid(t *testing.T) {
	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	var val *struct{}
	err := configure(val, v, flagSet, nil)

	assert.ErrorIs(t, err, ErrNotValidValue)
}

func TestConfigure_Primitives(t *testing.T) {
	type TestI interface {
		SomeFunc()
	}

	type root struct {
		Field1 string
		Field2 string `env:"-"`
		Field3 interface{}
		Field4 *interface{}
		Field5 string `env:",squash"`
		Field6 TestI
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	cfg := &root{}
	items := make(map[string]*Field)
	err := configure(cfg, v, flagSet, itemsToMapCallback(items))

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Contains(t, items, "field1")
	assert.Contains(t, items, "field5")
}

func TestConfigure_NamedStructure(t *testing.T) {
	type Inner struct {
		SubField1 string
	}

	type root struct {
		Field1 string
		Field2 Inner
		Field3 *Inner
		Field4 Inner `env:",squash"`
		Field5 Inner `env:"-"`
		Inner
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	cfg := &root{}
	items := make(map[string]*Field)
	err := configure(cfg, v, flagSet, itemsToMapCallback(items))

	assert.NoError(t, err)
	assert.Len(t, items, 9)
	assert.Contains(t, items, "field1")
	assert.Contains(t, items, "field2")
	assert.Contains(t, items, "field2.subField1")
	assert.Contains(t, items, "field3")
	assert.Contains(t, items, "field3.subField1")
	assert.Contains(t, items, "subField1")
	assert.Contains(t, items, "inner")
	assert.Contains(t, items, "inner.subField1")
}

func TestConfigure_StructureSquashDuplication1(t *testing.T) {
	type Inner struct {
		Field1 string
	}

	type root struct {
		Field1 string
		Field2 Inner `env:",squash"`
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	cfg := &root{}
	err := configure(cfg, v, flagSet, nil)

	assert.ErrorIs(t, err, ErrDuplicateKey)
}

func TestConfigure_StructureSquashDuplication2(t *testing.T) {
	type Inner struct {
		Field1 string
	}

	type root struct {
		Field1 Inner
		Field2 Inner `env:",squash"`
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	cfg := &root{}
	err := configure(cfg, v, flagSet, nil)

	assert.ErrorIs(t, err, ErrDuplicateKey)
}

func TestConfigure_Time(t *testing.T) {
	type root struct {
		Field1 time.Time
		Field2 time.Duration
		Field3 *time.Time
		Field4 *time.Duration
	}

	v := viper.New()
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	cfg := &root{}
	items := make(map[string]*Field)
	err := configure(cfg, v, flagSet, itemsToMapCallback(items))

	assert.NoError(t, err)
	assert.Len(t, items, 4)
	assert.Contains(t, items, "field1")
	assert.Contains(t, items, "field2")
	assert.Contains(t, items, "field3")
	assert.Contains(t, items, "field4")
}

func Test_StringToTimeRFC3339Decoder(t *testing.T) {
	str := time.Now().Format(time.RFC3339)
	m := map[string]interface{}{"field1": str, "field2": str}
	type root struct {
		Field1 time.Time
		Field2 *time.Time
	}

	hook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeHookFunc(time.RFC3339),
	)

	v := &root{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     v,
		DecodeHook: hook,
	})

	assert.NoError(t, err)

	err = dec.Decode(m)
	assert.NoError(t, err)
	assert.Equal(t, str, v.Field1.Format(time.RFC3339))
	assert.Equal(t, str, v.Field2.Format(time.RFC3339))
}
