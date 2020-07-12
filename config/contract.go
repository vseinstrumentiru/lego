package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagBinder interface {
	To(key string)
}

type Env interface {
	SetDefault(key string, value interface{})
	SetAlias(alias string, originalKey string)
	SetFlag(name string, value interface{}, usage string) FlagBinder
}

type Validatable interface {
	Validate() (err error)
}

type WithDefaults interface {
	SetDefaults(env Env)
}

// Deprecated: use WithDefaults interface
type WithDefaultsDeprecated interface {
	SetDefaults(key string, env *viper.Viper, flag *pflag.FlagSet)
}
