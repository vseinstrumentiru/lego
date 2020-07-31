package config

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
	"time"
)

type flagBinder struct {
	name  string
	value interface{}
	env   *env
}

func (f *flagBinder) To(key string) {
	_ = f.env.viper.BindPFlag(f.env.path+key, f.env.flag.Lookup(f.name))
	f.env.SetDefault(key, f.value)
}

type env struct {
	path  string
	viper *viper.Viper
	flag  *pflag.FlagSet
}

func (e *env) WithKey(key string) Env {
	if key == "" {
		return e
	}

	return &env{
		path:  e.path + key + ".",
		viper: e.viper,
		flag:  e.flag,
	}
}

func (e *env) SetDefault(key string, value interface{}) {
	e.viper.SetDefault(e.path+key, value)
}

func (e *env) SetAlias(alias string, originalKey string) {
	e.viper.RegisterAlias(alias, e.path+originalKey)
}

func (e *env) SetFlag(name string, value interface{}, usage string) FlagBinder {
	switch val := value.(type) {
	case bool:
		e.flag.Bool(name, val, usage)
	case []bool:
		e.flag.BoolSlice(name, val, usage)
	case []byte:
		e.flag.BytesHex(name, val, usage)
	case time.Duration:
		e.flag.Duration(name, val, usage)
	case []time.Duration:
		e.flag.DurationSlice(name, val, usage)
	case float32:
		e.flag.Float32(name, val, usage)
	case []float32:
		e.flag.Float32Slice(name, val, usage)
	case float64:
		e.flag.Float64(name, val, usage)
	case []float64:
		e.flag.Float64Slice(name, val, usage)
	case int:
		e.flag.Int(name, val, usage)
	case int8:
		e.flag.Int8(name, val, usage)
	case int16:
		e.flag.Int16(name, val, usage)
	case int32:
		e.flag.Int32(name, val, usage)
	case int64:
		e.flag.Int64(name, val, usage)
	case []int:
		e.flag.IntSlice(name, val, usage)
	case []int32:
		e.flag.Int32Slice(name, val, usage)
	case []int64:
		e.flag.Int64Slice(name, val, usage)
	case net.IP:
		e.flag.IP(name, val, usage)
	case []net.IP:
		e.flag.IPSlice(name, val, usage)
	case net.IPMask:
		e.flag.IPMask(name, val, usage)
	case net.IPNet:
		e.flag.IPNet(name, val, usage)
	default:
		emperror.Panic(errors.NewWithDetails("flag type not found", "name", name))
	}

	return &flagBinder{
		name:  name,
		value: value,
		env:   e,
	}
}
