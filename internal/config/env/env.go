package env

import (
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/vseinstrumentiru/lego/config"
)

type Env interface {
	config.Env
	RootKey() string
	Sub(key string) config.Env
	Load(config interface{}, keys ...string) error
	OnFlag(name string, callback interface{}) bool
	Get(key string) interface{}
}

type ErrConfigFileNotFound = viper.ConfigFileNotFoundError

func IsFileNotFound(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)

	return ok
}

func New(path string) Env {
	v, p := viper.New(), pflag.NewFlagSet("lego", pflag.ExitOnError)
	v.SetEnvPrefix(path)

	return &env{
		viper: v,
		flag:  p,
	}
}

type env struct {
	path  string
	viper *viper.Viper
	flag  *pflag.FlagSet
}

func (e *env) Get(key string) interface{} {
	return e.viper.Get(key)
}

func (e *env) OnFlag(name string, callback interface{}) bool {
	val := reflect.ValueOf(callback)

	if val.Kind() != reflect.Func {
		emperror.Panic(errors.New("callback must be a function"))
	}

	if val.Type().NumIn() != 1 {
		emperror.Panic(errors.New("callback must accept only one argument"))
	}

	if val.Type().NumIn() != 1 {
		emperror.Panic(errors.New("callback must accept only one argument"))
	}

	if val.Type().In(0).Kind() == reflect.Ptr {
		emperror.Panic(errors.New("callback must argument must be non-pointer value"))
	}

	t := reflect.New(val.Type().In(0)).Elem().Interface()
	var arg interface{}

	switch t.(type) {
	case string:
		if v, err := e.flag.GetString(name); err == nil && v != "" {
			arg = v
		}
	case []string:
		if v, err := e.flag.GetStringSlice(name); err == nil && len(v) != 0 {
			arg = v
		}
	case bool:
		if v, err := e.flag.GetBool(name); err == nil && v {
			arg = v
		}
	case []bool:
		if v, err := e.flag.GetBoolSlice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case []byte:
		if v, err := e.flag.GetBytesHex(name); err == nil && len(v) > 0 {
			arg = v
		}
	case time.Duration:
		if v, err := e.flag.GetDuration(name); err == nil && v > 0 {
			arg = v
		}
	case []time.Duration:
		if v, err := e.flag.GetDurationSlice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case float32:
		if v, err := e.flag.GetFloat32(name); err == nil {
			arg = v
		}
	case []float32:
		if v, err := e.flag.GetFloat32Slice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case float64:
		if v, err := e.flag.GetFloat64(name); err == nil {
			arg = v
		}
	case []float64:
		if v, err := e.flag.GetFloat64Slice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case int:
		if v, err := e.flag.GetInt(name); err == nil {
			arg = v
		}
	case int8:
		if v, err := e.flag.GetInt8(name); err == nil {
			arg = v
		}
	case int16:
		if v, err := e.flag.GetInt16(name); err == nil {
			arg = v
		}
	case int32:
		if v, err := e.flag.GetInt32(name); err == nil {
			arg = v
		}
	case int64:
		if v, err := e.flag.GetInt64(name); err == nil {
			arg = v
		}
	case []int:
		if v, err := e.flag.GetIntSlice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case []int32:
		if v, err := e.flag.GetInt32Slice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case []int64:
		if v, err := e.flag.GetInt64Slice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case net.IP:
		if v, err := e.flag.GetIP(name); err == nil && len(v) > 0 {
			arg = v
		}
	case []net.IP:
		if v, err := e.flag.GetIPSlice(name); err == nil && len(v) > 0 {
			arg = v
		}
	case net.IPMask:
		if v, err := e.flag.GetIPv4Mask(name); err == nil && len(v) > 0 {
			arg = v
		}
	case net.IPNet:
		if v, err := e.flag.GetIPNet(name); err == nil && len(v.IP) > 0 {
			arg = v
		}
	default:
		emperror.Panic(errors.NewWithDetails("flag type not found", "name", name))
	}

	if arg == nil {
		return false
	}

	val.Call([]reflect.Value{reflect.ValueOf(arg)})
	return true
}

func (e *env) RootKey() string {
	return e.path
}

func (e *env) Sub(key string) config.Env {
	if key == "" || key == "." {
		return e
	}

	key = strings.TrimRight(key, ".")

	return &env{
		path:  key,
		viper: e.viper,
		flag:  e.flag,
	}
}

func (e *env) SetDefault(key string, value interface{}) {
	e.viper.SetDefault(e.path+"."+key, value)
}

func (e *env) SetAlias(alias string, originalKey string) {
	e.viper.RegisterAlias(alias, e.path+"."+originalKey)
}

func (e *env) SetFlag(name string, value interface{}, usage string) config.FlagBinder {
	switch val := value.(type) {
	case string:
		e.flag.String(name, val, usage)
	case []string:
		e.flag.StringSlice(name, val, usage)
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

func (e *env) Load(config interface{}, keys ...string) error {
	e.viper.AddConfigPath(".")
	e.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	e.viper.AllowEmptyEnv(true)
	e.viper.AutomaticEnv()

	for i := 0; i < len(keys); i++ {
		err := e.viper.BindEnv(keys[i], strings.ToUpper(strings.ReplaceAll(e.path+"."+keys[i], ".", "_")))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	e.flag.String("config", "", "Configuration file")
	e.flag.String("config-path", "", "Search path for configuration file")

	_ = e.flag.Parse(os.Args[1:])

	e.OnFlag("config", e.viper.SetConfigFile)
	e.OnFlag("config-path", e.viper.AddConfigPath)

	if err := e.viper.ReadInConfig(); err != nil && !IsFileNotFound(err) {
		return errors.Wrap(err, "failed to read configuration")
	}

	hook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		StringToTypeDecoder,
	)

	if err := e.viper.Unmarshal(config, viper.DecodeHook(hook)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type flagBinder struct {
	name  string
	value interface{}
	env   *env
}

func (f *flagBinder) To(key string) {
	_ = f.env.viper.BindPFlag(f.env.path+key, f.env.flag.Lookup(f.name))
	f.env.SetDefault(key, f.value)
}
