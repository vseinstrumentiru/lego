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

func NewDefaultConfig() *Config {
	return &Config{Port: 8080, ShutdownTimeout: 5 * time.Second}
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("port", 8080)
	env.SetDefault("shutdownTimeout", "5s")
}
