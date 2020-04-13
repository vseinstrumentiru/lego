package lego

import (
	"net"
	"time"
)

type Process interface {
	LogErr
	Name() string
	Listen(network, addr string) (net.Listener, error)
	Run(execute func() error, interrupt func(error))
	ShutdownTimeout() time.Duration
}
