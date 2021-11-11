package config

import (
	"os"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Parser is a config parser interface
type Parser interface {
	// Unmarshal reads values from file and/or environment,
	// sets them to cfg and validates cfg
	Unmarshal(cfg interface{}) error
}

// Option is a config parser constructor's option
type Option func(p *parser)

// New is a config parser constructor
func New(opts ...Option) *parser {
	p := &parser{
		viper: viper.New(),
		flag:  pflag.NewFlagSet(os.Args[0], pflag.ExitOnError),
		decoders: []mapstructure.DecodeHookFunc{
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			StringToTypeDecoder,
		},
		args: os.Args[1:],
	}

	p.flag.ParseErrorsWhitelist.UnknownFlags = true

	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}

	return p
}

type parser struct {
	viper    *viper.Viper
	flag     *pflag.FlagSet
	decoders []mapstructure.DecodeHookFunc
	args     []string
}

// Unmarshal reads values from file and/or environment,
// sets them to cfg and validates cfg
func (p *parser) Unmarshal(cfg interface{}) error {
	if err := configure(cfg, p.viper, p.flag, nil); err != nil {
		return err
	}

	p.viper.AllowEmptyEnv(true)
	p.viper.AutomaticEnv()
	p.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	p.viper.AddConfigPath(".")

	if err := p.viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return errors.Wrap(err, "can't read configuration")
	}

	mapStructureConfig := func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(p.decoders...)
		config.TagName = "env"
	}

	if err := p.viper.Unmarshal(cfg, mapStructureConfig); err != nil {
		return errors.Wrap(err, "can't unmarshal config")
	}

	if err := p.flag.Parse(p.args); err != nil {
		return errors.Wrap(err, "can't parse args")
	}

	err := p.validate(cfg)

	return err
}

func (p *parser) validate(cfg interface{}) error {
	err := validator.New().Struct(cfg)

	if err != nil {
		validatorErr := err.(validator.ValidationErrors)
		allErr := make([]error, 0, len(validatorErr))
		for _, e := range validatorErr {
			allErr = append(allErr, e)
		}

		err = errors.Combine(allErr...)
	}

	return err
}
