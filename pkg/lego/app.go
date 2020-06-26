package lego

import (
	"github.com/vseinstrumentiru/lego/internal/lego"
)

type App = lego.App

type AppWithConfig interface {
	lego.AppWithConfig
}

type AppWithHttp interface {
	lego.AppWithHttp
}

type AppWithGrpc interface {
	lego.AppWithGrpc
}

type AppWithEventHandlers interface {
	lego.AppWithEventHandlers
}

type AppWithPublishers interface {
	lego.AppWithPublishers
}

type AppWithStats interface {
	lego.AppWithStats
}

type AppWithRegistration interface {
	lego.AppWithRegistration
}

type AppWithRunner interface {
	lego.AppWithRunner
}
