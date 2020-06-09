package lego

import (
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
)

type WithSwitch struct {
	Enabled bool
}

type WithCustomConfig interface {
	GetConfig() lego2.Config
	SetConfig(config lego2.Config)
}
