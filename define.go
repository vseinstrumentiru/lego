package LeGo

import (
	"github.com/vseinstrumentiru/lego/pkg/lego"
)

type App = lego.App

func extractConfig(target interface{}) lego.Config {
	if cTarget, ok := target.(lego.AppWithConfig); ok {
		return cTarget.GetConfig()
	}

	return nil
}

func provideConfig(target interface{}, cfg lego.Config) {
	if cTarget, ok := target.(lego.AppWithConfig); ok {
		cTarget.SetConfig(cfg)
	}
}
