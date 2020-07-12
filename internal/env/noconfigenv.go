package env

import (
	base "github.com/vseinstrumentiru/lego/config"
)

func NewNoConfigEnv(path string) Env {
	return &noConfigEnv{NewBaseEnv(path)}
}

type noConfigEnv struct {
	*baseEnv
}

func (e *noConfigEnv) Configure() error {
	if err := e.setEnv([]string{"name", "dataCenter"}); err != nil {
		return err
	}

	app := base.Undefined()

	if name := e.viper.GetString("name"); name != "" {
		app = &base.Application{
			Name:       name,
			DataCenter: e.viper.GetString("dataCenter"),
		}
	}

	e.instances = append(e.instances, instance{val: app})

	return nil
}
