package http

import (
	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	Port     int
	IsPublic bool
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("port", "8081")
}
