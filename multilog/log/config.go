package log

import (
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	Color bool
	Stop  bool
	Depth int
	Level logur.Level
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("depth", -1)
}
