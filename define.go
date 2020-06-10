package LeGo

import (
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
)

type App = lego2.App

func extractConfig(target interface{}) lego2.Config {
	if cTarget, ok := target.(lego2.AppWithConfig); ok {
		return cTarget.GetConfig()
	}

	return nil
}

func provideConfig(target interface{}, cfg lego2.Config) {
	if cTarget, ok := target.(lego2.AppWithConfig); ok {
		cTarget.SetConfig(cfg)
	}
}
