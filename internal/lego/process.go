package lego

import (
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/vseinstrumentiru/lego/internal/lego/build"
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
	RegisterCheck(cfg *health.Config) error
	ShutdownTimeout() time.Duration
}
