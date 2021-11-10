package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// WithEnvPrefix adds prefix to environment variables
// For prefix=myapp all variables must starts with MYAPP_
// see https://github.com/spf13/viper
func WithEnvPrefix(prefix string) Option {
	return func(p *parser) {
		p.viper.SetEnvPrefix(prefix)
	}
}

// WithFlagSet sets custom pflag.FlagSet
func WithFlagSet(set *pflag.FlagSet) Option {
	return func(p *parser) {
		p.flag = set
	}
}

// WithViper sets custom viper.Viper
func WithViper(v *viper.Viper) Option {
	return func(p *parser) {
		p.viper = v
	}
}

// WithDecodeFunctions adds custom decode function
// for decoding string into your type
func WithDecodeFunctions(decoders ...mapstructure.DecodeHookFunc) Option {
	return func(p *parser) {
		p.decoders = append(p.decoders, decoders...)
	}
}

// WithArgs sets args for parsing, useful for tests
// default: os.Args[1:]
func WithArgs(args []string) Option {
	return func(p *parser) {
		p.args = args
	}
}
