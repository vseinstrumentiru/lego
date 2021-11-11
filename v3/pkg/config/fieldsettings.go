package config

import (
	"net"
	"time"

	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func bindEnv(field *Field, v *viper.Viper, _ *pflag.FlagSet) error {
	if field.isSquashed || field.hasFields {
		return nil
	}

	return v.BindEnv(field.key())
}

func setDefaults(field *Field, _ *viper.Viper, _ *pflag.FlagSet) error {
	withDefaults := field.value.Interface().(WithDefaults)
	withDefaults.SetDefaults()

	return nil
}

func setFlag(name, short, usage string) fieldSetting {
	return func(field *Field, _ *viper.Viper, flagSet *pflag.FlagSet) error {

		i := field.value.Interface()

		switch val := i.(type) {
		case *string:
			flagSet.StringVarP(val, name, short, *val, usage)
		case *[]string:
			flagSet.StringSliceVarP(val, name, short, *val, usage)
		case *bool:
			flagSet.BoolVarP(val, name, short, *val, usage)
		case *[]bool:
			flagSet.BoolSliceVarP(val, name, short, *val, usage)
		case *[]byte:
			flagSet.BytesHexVarP(val, name, short, *val, usage)
		case *time.Duration:
			flagSet.DurationVarP(val, name, short, *val, usage)
		case *[]time.Duration:
			flagSet.DurationSliceVarP(val, name, short, *val, usage)
		case *float32:
			flagSet.Float32VarP(val, name, short, *val, usage)
		case *[]float32:
			flagSet.Float32SliceVarP(val, name, short, *val, usage)
		case *float64:
			flagSet.Float64VarP(val, name, short, *val, usage)
		case *[]float64:
			flagSet.Float64SliceVarP(val, name, short, *val, usage)
		case *int:
			flagSet.IntVarP(val, name, short, *val, usage)
		case *int8:
			flagSet.Int8VarP(val, name, short, *val, usage)
		case *int16:
			flagSet.Int16VarP(val, name, short, *val, usage)
		case *int32:
			flagSet.Int32VarP(val, name, short, *val, usage)
		case *int64:
			flagSet.Int64VarP(val, name, short, *val, usage)
		case *[]int:
			flagSet.IntSliceVarP(val, name, short, *val, usage)
		case *[]int32:
			flagSet.Int32SliceVarP(val, name, short, *val, usage)
		case *[]int64:
			flagSet.Int64SliceVarP(val, name, short, *val, usage)
		case *net.IP:
			flagSet.IPVarP(val, name, short, *val, usage)
		case *[]net.IP:
			flagSet.IPSliceVarP(val, name, short, *val, usage)
		case *net.IPMask:
			flagSet.IPMaskVarP(val, name, short, *val, usage)
		case *net.IPNet:
			flagSet.IPNetVarP(val, name, short, *val, usage)
		default:
			return errors.NewWithDetails("flag type not found", "name", name)
		}

		return nil
	}
}
