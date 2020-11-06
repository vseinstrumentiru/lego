package env

import (
	"net"
	"os"
	"strings"
	"time"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/vseinstrumentiru/lego/common/cast"
	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/inject"
	di "github.com/vseinstrumentiru/lego/internal/container"
)

func NewBaseEnv(path string) *baseEnv {
	v, p := viper.New(), pflag.NewFlagSet("lego", pflag.ExitOnError)
	v.SetEnvPrefix(path)

	return &baseEnv{
		viper: v,
		flag:  p,
	}
}

type instance struct {
	val interface{}
	key string
}

type baseEnv struct {
	path      string
	viper     *viper.Viper
	flag      *pflag.FlagSet
	instances []instance
}

func (e *baseEnv) ConfigureInstances(container di.Container) error {
	for _, i := range e.instances {
		var opts []inject.RegisterOption
		if i.key != "" {
			opts = append(opts, di.WithName(i.key))
		}

		if err := container.Instance(i.val, opts...); err != nil {
			return err
		}
	}

	return nil
}

func (e *baseEnv) Get(key string) interface{} {
	return e.viper.Get(key)
}

func (e *baseEnv) OnFlag(name string, callback interface{}) bool {
	return cast.OnCheck(e.flag, name, callback)
}

func (e *baseEnv) RootKey() string {
	return e.path
}

func (e *baseEnv) Sub(key string) config.Env {
	if key == "" || key == "." {
		return e
	}

	key = strings.TrimRight(key, ".")

	return &baseEnv{
		path:  key,
		viper: e.viper,
		flag:  e.flag,
	}
}

func (e *baseEnv) SetDefault(key string, value interface{}) {
	path := ""
	if e.path != "" {
		path = e.path + "."
	}
	e.viper.SetDefault(path+key, value)
}

func (e *baseEnv) SetAlias(alias string, originalKey string) {
	e.viper.RegisterAlias(alias, e.path+"."+originalKey)
}

func (e *baseEnv) SetFlag(name string, value interface{}, usage string) config.FlagBinder {
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

func (e *baseEnv) setEnv(keys []string) error {
	e.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	e.viper.AllowEmptyEnv(true)
	e.viper.AutomaticEnv()

	for i := 0; i < len(keys); i++ {
		err := e.viper.BindEnv(keys[i], strings.ToUpper(strings.ReplaceAll(e.path+"."+keys[i], ".", "_")))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	_ = e.flag.Parse(os.Args[1:])

	return nil
}

type flagBinder struct {
	name  string
	value interface{}
	env   *baseEnv
}

func (f *flagBinder) To(key string) {
	_ = f.env.viper.BindPFlag(f.env.path+key, f.env.flag.Lookup(f.name))
	f.env.SetDefault(key, f.value)
}
