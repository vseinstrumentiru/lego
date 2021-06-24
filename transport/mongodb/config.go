package mongodb

import (
	"emperror.dev/errors"
	"github.com/lebrains/gomongowrapper"
	"github.com/vseinstrumentiru/lego/v2/config"
	"time"
)

var ErrEmptyName = errors.New("database name must be set")

type Config struct {
	gomongowrapper.Config `mapstructure:",squash"  load:"true"`
	ConnectTimeout        time.Duration `optional:"true"`
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("connectTimeout", 5*time.Second)
}

func (c Config) Validate() (err error) {
	if e := c.Config.Validate(); e != nil {
		err = errors.Append(err, e)
	}

	if c.Name == "" {
		err = errors.Append(err, ErrEmptyName)
	}

	return
}
