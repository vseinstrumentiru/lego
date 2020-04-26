package config

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"fmt"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/pkg/build"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"os"
	"reflect"
	"strings"
)

const (
	defaultEnvPrefix = "app"
)

type ErrConfigFileNotFound = viper.ConfigFileNotFoundError

func IsFileNotFound(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)

	return ok
}

func configure() (*viper.Viper, *pflag.FlagSet) {
	v, p := viper.New(), pflag.NewFlagSet("lego", pflag.ExitOnError)

	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	p.String("config", "", "Configuration file")
	p.String("config-path", "", "Search path for configuration file")
	p.Bool("version", false, "Show version information")

	v.AddConfigPath(".")

	return v, p
}

func prepareCustomConfig(customCfg lego.Config) (lego.Config, bool) {
	if customCfg == lego.Config(nil) {
		return nil, false
	}

	cfgVal := reflect.ValueOf(customCfg)
	if !cfgVal.IsValid() {
		emperror.Panic(errors.New("config is nil pointer"))
	}

	return customCfg, true
}

func buildConfig(v *viper.Viper, customCfg lego.Config, customCfgPrefix string) (cfg Config, err error) {
	builder := dynamicstruct.NewStruct().
		AddField("Srv", Config{}, "")

	var hasCustomCfg bool

	if customCfg != nil {
		hasCustomCfg = true
		envPrefix := customCfgPrefix

		if cApp, ok := customCfg.(lego.ConfigWithCustomEnvPrefix); ok {
			envPrefix = cApp.GetEnvPrefix()
		}
		v.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(envPrefix)))

		builder = builder.AddField("Custom", customCfg, fmt.Sprintf(`mapstructure:"%s"`, customCfgPrefix))
	}

	unmarshalStruct := builder.Build().New()
	if err = v.Unmarshal(&unmarshalStruct); err != nil {
		return Config{}, err
	}

	structReflect := reflect.ValueOf(unmarshalStruct).Elem()
	cfg = structReflect.FieldByName("Srv").Interface().(Config)
	if hasCustomCfg {
		cfg.Custom = structReflect.FieldByName("Custom").Interface().(lego.Config)
	}

	cfg.Build = build.New()

	return cfg, nil
}

type Option func(env *viper.Viper, flags *pflag.FlagSet)

func setServerCfgDefaults() Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		(Config{}).SetDefaults(env, flags)
	}
}

func setCustomCfgDefaults(customCfg lego.Config) Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		customCfg.SetDefaults(env, flags)
	}
}

func WithDefaultName(name string) Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		env.SetDefault("srv.name", name)
	}
}

func Provide(customCfg lego.Config, options ...Option) (Config, error) {
	v, p := configure()
	var hasAppCfg bool
	customCfg, hasAppCfg = prepareCustomConfig(customCfg)

	defaultOptions := []Option{setServerCfgDefaults()}
	if hasAppCfg {
		defaultOptions = append(defaultOptions, setCustomCfgDefaults(customCfg))
	}

	options = append(defaultOptions, options...)

	for _, opt := range options {
		opt(v, p)
	}

	_ = p.Parse(os.Args[1:])

	if c, _ := p.GetString("config"); c != "" {
		v.SetConfigFile(c)
	} else if c, _ := p.GetString("config-path"); c != "" {
		v.AddConfigPath(c)
	}

	returnErr := v.ReadInConfig()
	if !IsFileNotFound(returnErr) {
		emperror.Panic(errors.Wrap(returnErr, "failed to read configuration"))
	}

	cfg, err := buildConfig(v, customCfg, defaultEnvPrefix)
	emperror.Panic(errors.Wrap(err, "failed to unmarshal application configuration"))

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", cfg.Name, cfg.Build.Version, cfg.Build.CommitHash, cfg.Build.BuildDate)

		os.Exit(0)
	}

	cfg.Name = strings.ReplaceAll(cfg.Name, "-", "_")

	return cfg, returnErr
}
