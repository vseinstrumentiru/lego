package env

import (
	"github.com/vseinstrumentiru/lego/v2/config"
)

type Keys []string

type Config interface{}

type Env interface {
	config.Env
	RootKey() string
	Sub(key string) config.Env
	OnFlag(name string, callback interface{}) bool
	Get(key string) interface{}
}
