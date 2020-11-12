package env

import (
	"reflect"

	"emperror.dev/errors"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	base "github.com/vseinstrumentiru/lego/v2/config"
)

func NewConfigEnv(path string) Env {
	return &configEnv{NewBaseEnv(path)}
}

type configEnv struct {
	*baseEnv
}

func (e *configEnv) Configure(cfg Config) error {
	parsed, err := e.configureParser(cfg)

	if err != nil {
		return err
	}

	e.flag.String("config", "", "Configuration file")
	e.flag.String("config-path", "", "Search path for configuration file")

	e.setDefaults(parsed.defaults)

	if err = e.setEnv(parsed.keys); err != nil {
		return err
	}

	if err = e.loadConfig(cfg); err != nil {
		return err
	}

	if err = e.validate(parsed.validates); err != nil {
		return err
	}

	e.setInstances(parsed)

	if !parsed.exist(base.Application{}) {
		app := base.Undefined()

		if name := e.viper.GetString("name"); name != "" {
			app = &base.Application{
				Name:       name,
				DataCenter: e.viper.GetString("dataCenter"),
			}
		}

		e.instances = append(e.instances, Instance{Val: app, IsDefault: true})
	}

	return nil
}

func (e *configEnv) configureParser(cfg Config) (*parser, error) {
	parsed := newParser()

	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("config must be a pointer")
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type()))
	}

	err := parsed.scan(v, "", flags{})

	return parsed, err
}

func (e *configEnv) setDefaults(def []defaults) {
	for i := len(def) - 1; i >= 0; i-- {
		d := def[i]
		d.val.SetDefaults(e.Sub(d.key))
	}
}

func (e *configEnv) validate(validates map[string]base.Validatable) (err error) {
	for _, i := range validates {
		err = errors.Append(err, i.Validate())
	}

	return err
}

func (e *configEnv) setInstances(parsed *parser) {
	for _, coll := range parsed.configs {
		coll.Items[coll.DefaultKey].IsDefault = true
		for _, i := range coll.Items {
			e.instances = append(e.instances, i)
		}
	}
}

func (e *configEnv) loadConfig(cfg Config) error {
	e.viper.AddConfigPath(".")

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

	if err := e.viper.Unmarshal(cfg, viper.DecodeHook(hook)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
