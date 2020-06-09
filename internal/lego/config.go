package lego

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config interface {
	Validatable
	SetDefaults(env *viper.Viper, flag *pflag.FlagSet)
}

type ConfigWithKey interface {
	Validatable
	SetDefaults(key string, env *viper.Viper, flag *pflag.FlagSet)
}

type ConfigWithCustomEnvPrefix interface {
	Config
	GetEnvPrefix() string
}
