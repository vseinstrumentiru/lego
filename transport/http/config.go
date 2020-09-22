package http

import (
	"time"

	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	Port            int
	IsPublic        bool
	ShutdownTimeout time.Duration
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("port", 8080)
	env.SetDefault("shutdownTimeout", "5s")
}
