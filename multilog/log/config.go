package log

import (
	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	Color bool
	Stop  bool
	Depth int
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("depth", 7)
}
