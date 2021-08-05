package metrics

import "github.com/vseinstrumentiru/lego/v2/config"

type Config struct {
	Port  int
	Debug bool
}

func (c *Config) SetDefaults(env config.Env) {
	env.SetFlag("telemetry-addr", 10000, "Telemetry HTTP server address").To("port")
}

func NewDefaultConfig() *Config {
	return &Config{Port: 10000}
}
