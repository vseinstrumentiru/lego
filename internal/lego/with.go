package lego

type WithSwitch struct {
	Enabled bool
}

type WithCustomConfig interface {
	GetConfig() Config
	SetConfig(config Config)
}
