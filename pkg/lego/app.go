package lego

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorilla/mux"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"io"
)

type App interface {
	GetName() string
	SetLogErr(logErr LogErr)
}

type AppWithConfig interface {
	GetConfig() Config
	SetConfig(config Config)
}

type AppWithHttp interface {
	RegisterHTTP(router *mux.Router) error
}

type AppWithGrpc interface {
	RegisterGRPC(server *grpc.Server) error
}

type AppWithSubscribers interface {
	RegisterEventHandlers(router *message.Router, subscriber message.Subscriber) error
}

type AppWithPublishers interface {
	RegisterEventDispatcher(publisher message.Publisher) error
}

type AppWithHealthChecks interface {
	GetHealthChecks() []*view.View
}

type AppWithRegistration interface {
	Register(p Process) (io.Closer, error)
}
