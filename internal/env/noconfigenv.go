package env

import (
	base "github.com/vseinstrumentiru/lego/v2/config"
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

	app := base.UndefinedApplication()

	if name := e.viper.GetString("name"); name != "" {
		app = &base.Application{
			Name:       name,
			DataCenter: e.viper.GetString("dataCenter"),
		}
	}

	e.instances = append(e.instances, Instance{Val: app, IsDefault: true})

	return nil
}
