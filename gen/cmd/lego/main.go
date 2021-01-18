package main

import (
	"github.com/gobuffalo/packr/v2"

	"github.com/vseinstrumentiru/lego/v2/app"
	"github.com/vseinstrumentiru/lego/v2/di"
	"github.com/vseinstrumentiru/lego/v2/gen/internal/newproject"
)

func main() {
	box := packr.New("assets", "../../assets")

	app.NewRuntime(
		app.EnvPath("lego"),
		app.NoDefaultProviders(),
		app.Provide(
			di.ProvideCommand(newproject.NewCommand(box)),
		),
	).Run()
}
