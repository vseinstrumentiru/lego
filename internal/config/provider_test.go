package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type appConfig struct {
	Some struct {
		Name string
		Any  string
	}
}

func (a appConfig) Validate() error {
	return nil
}

func (a appConfig) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.Set("app.some.name", "Asd")
	env.Set("app.some.any", "Test")
	env.Set("app.some.name", "Test")
	env.SetDefault("srv.name", "TestSrv")
}

func Test_Provider(t *testing.T) {
	ass := assert.New(t)

	var cfg appConfig
	srv, err := Provide(cfg)

	ass.IsType(err, viper.ConfigFileNotFoundError{})

	cfg = srv.App.(appConfig)

	ass.Equal("TestSrv", srv.Name)
	ass.Equal("Test", cfg.Some.Name)
}
