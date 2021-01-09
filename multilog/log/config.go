package log

import (
	"time"

	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
)

func DefaultConfig() *Config {
	return &Config{
		Color:      false,
		Stop:       false,
		Depth:      -1,
		Level:      logur.Trace,
		TimeFormat: "15:04:05.000",
	}
}

type Config struct {
	Color      bool
	Stop       bool
	Depth      int
	Level      logur.Level
	TimeFormat string
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("depth", -1)
	env.SetDefault("timeFormat", time.RFC3339Nano)
}
