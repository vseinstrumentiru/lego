package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type app struct {
	Name string
	Some struct {
		Val1 string
		Val2 string
	}
}

func (a app) Validate() error { return nil }
func (a app) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	setDefaults(env, flag)
}

func setDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("app.some.val2", "val2-default")
	env.SetDefault("srv.name", "test-default")
	env.SetDefault("srv.events.nats.clientid", "test")

	env.AddConfigPath("./test")
	env.SetConfigName("test")
}

func Test_ConfigDefaultsWithoutApp(t *testing.T) {
	ass := assert.New(t)

	srv, err := Provide(nil, setDefaults)

	ass.Nil(err)
	ass.Equal("test", srv.Name)
	ass.Equal("nats://testaddr", srv.Events.Nats.Addr)
	ass.Equal("test", srv.Events.Nats.ClientID)
}

func Test_ConfigDefaultsWithNillPointerApp(t *testing.T) {
	ass := assert.New(t)
	var a *app
	ass.Panics(func() {
		Provide(a)
	})
}

func Test_ConfigDefaultsWithPointerApp(t *testing.T) {
	ass := assert.New(t)
	srv, err := Provide(&app{})

	ass.Nil(err)
	ass.Equal("test", srv.Name)
	ass.Equal("nats://testaddr", srv.Events.Nats.Addr)
	ass.Equal("test", srv.Events.Nats.ClientID)

	ass.NotNil(srv.App)
	ass.NotNil("app-test", srv.App.(*app).Name)
}

func Test_ConfigDefaultsWithApp(t *testing.T) {
	ass := assert.New(t)
	var a app
	srv, err := Provide(a)

	ass.Nil(err)

	ass.Equal("test", srv.Name)
	ass.Equal("nats://testaddr", srv.Events.Nats.Addr)
	ass.Equal("test", srv.Events.Nats.ClientID)

	ass.NotNil(srv.App)
	ass.NotNil("app-test", srv.App.(app).Name)
}
