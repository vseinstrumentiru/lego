package main

import (
	"github.com/gobuffalo/packr/v2"

	"github.com/vseinstrumentiru/lego/v2/app"
	"github.com/vseinstrumentiru/lego/v2/di"
)

var box = packr.New("assets", "../../assets")

func main() {
	app.NewRuntime(
		app.EnvPath("lego"),
		app.NoDefaultProviders(),
		app.Provide(
			di.ProvideCommand(newProjectCMD),
		),
	).Run()
}
