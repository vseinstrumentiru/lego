package main

import (
	"github.com/vseinstrumentiru/lego/v2/app"
	"github.com/vseinstrumentiru/lego/v2/di"
	"github.com/vseinstrumentiru/lego/v2/gen/internal/app/project"
	"github.com/vseinstrumentiru/lego/v2/gen/internal/app/service"
	"github.com/vseinstrumentiru/lego/v2/gen/internal/app/store"
)

func main() {
	app.NewRuntime(
		app.EnvPath("lego"),
		app.Provide(
			di.ProvideCommand(project.NewCommand()),
			di.ProvideCommand(service.NewCommand()),
			di.ProvideCommand(store.NewCommand()),
		),
	).Run()
}
