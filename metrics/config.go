package metrics

type Config struct {
	Port  int
	Debug bool
}

func NewDefaultConfig() *Config {
	return &Config{Port: 10000}
}
