package grpc

import (
	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	Port     int
	IsPublic bool
}

func NewDefaultConfig() *Config {
	return &Config{Port: 8081}
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("port", "8081")
}
