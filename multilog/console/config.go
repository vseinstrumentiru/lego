package console

import (
	"time"

	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
)

const (
	JSONFormat  = "json"
	HumanFormat = "human"
)

const (
	DefaultTimeFormat = "15:04:05.000"
)

func DefaultConfig() *Config {
	return &Config{
		Format:     JSONFormat,
		Stop:       false,
		Depth:      -1,
		Level:      logur.Trace,
		TimeFormat: DefaultTimeFormat,
	}
}

type Config struct {
	// Deprecated: use Format: human
	Color      bool
	Format     string
	Stop       bool
	Depth      int
	Level      logur.Level
	TimeFormat string
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("depth", -1)
	env.SetDefault("format", JSONFormat)
	env.SetDefault("timeFormat", time.RFC3339Nano)
}
