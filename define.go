package LeGo

import (
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/pkg/lego"
)

type App = lego.App

func extractConfig(target interface{}) lego2.Config {
	if cTarget, ok := target.(lego.AppWithConfig); ok {
		return cTarget.GetConfig()
	}

	return nil
}

func provideConfig(target interface{}, cfg lego2.Config) {
	if cTarget, ok := target.(lego.AppWithConfig); ok {
		cTarget.SetConfig(cfg)
	}
}
