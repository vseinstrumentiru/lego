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
	if c, _ := p.GetString("config"); c != "" {
		v.SetConfigFile(c)
	} else if c, _ := p.GetString("config-path"); c != "" {
		v.AddConfigPath(c)
	}

	return v, p
}

func prepareAppConfig(cfg lego.Config) (lego.Config, bool) {
	if cfg == lego.Config(nil) {
		return nil, false
	}

	cfgVal := reflect.ValueOf(cfg)
	if !cfgVal.IsValid() {
		emperror.Panic(errors.New("app config is nil pointer"))
	}

	return cfg, true
}

func buildConfig(v *viper.Viper, app lego.Config) (srv Server, err error) {
	builder := dynamicstruct.NewStruct().
		AddField("Srv", Server{}, "")

	var hasApp bool

	if app != nil {
		hasApp = true
		envPrefix := defaultEnvPrefix

		if cApp, ok := app.(lego.ConfigWithCustomEnvPrefix); ok {
			envPrefix = cApp.GetEnvPrefix()
		}
		v.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(envPrefix)))

		builder = builder.AddField("App", app, fmt.Sprintf(`mapstructure:"%s"`, envPrefix))
	}

	unmarshalStruct := builder.Build().New()
	if err = v.Unmarshal(&unmarshalStruct); err != nil {
		return Server{}, err
	}

	structReflect := reflect.ValueOf(unmarshalStruct).Elem()
	srv = structReflect.FieldByName("Srv").Interface().(Server)
	if hasApp {
		srv.App = structReflect.FieldByName("App").Interface().(lego.Config)
	}

	srv.Build = build.New()

	return srv, nil
}

type Option func(env *viper.Viper, flags *pflag.FlagSet)

func setServerDefaults() Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		(Server{}).SetDefaults(env, flags)
	}
}

func setAppDefaults(app lego.Config) Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		app.SetDefaults(env, flags)
	}
}

func WithDefaultName(name string) Option {
	return func(env *viper.Viper, flags *pflag.FlagSet) {
		env.SetDefault("srv.name", name)
	}
}

func Provide(app lego.Config, options ...Option) (Server, error) {
	v, p := configure()
	var hasAppCfg bool
	app, hasAppCfg = prepareAppConfig(app)

	defaultOptions := []Option{setServerDefaults()}
	if hasAppCfg {
		defaultOptions = append(defaultOptions, setAppDefaults(app))
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

	server, err := buildConfig(v, app)
	emperror.Panic(errors.Wrap(err, "failed to unmarshal application configuration"))

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", server.Name, server.Build.Version, server.Build.CommitHash, server.Build.BuildDate)

		os.Exit(0)
	}

	return server, returnErr
}
