package multilog

import (
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
)

type Config struct {
	Level logur.Level
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("level", logur.Error)
}
