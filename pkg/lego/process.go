package lego

import (
	"github.com/vseinstrumentiru/lego/pkg/build"
	"net"
	"time"
)

type Process interface {
	LogErr

	Name() string
	DataCenterName() string
	Build() build.Info
	Env() string
	IsDebug() bool

	Listen(network, addr string) (net.Listener, error)
	Background(execute func() error, interrupt func(error))
	ShutdownTimeout() time.Duration
}
