package http

import (
	"time"

	"github.com/vseinstrumentiru/lego/v2/config"
)

type Config struct {
	Port            int
	IsPublic        bool
	ShutdownTimeout time.Duration
}

func NewDefaultConfig() *Config {
	return &Config{Port: 8080, ShutdownTimeout: 5 * time.Second, IsPublic: false}
}

func (c *Config) SetDefaults(env config.Env) {
	env.SetDefault("port", 8080)
	env.SetDefault("shutdownTimeout", "5s")
	env.SetFlag("http-addr", 8080, "App HTTP server address").To("port")
}
