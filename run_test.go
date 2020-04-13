package LeGo

import (
	"context"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"io"
	"testing"
	"time"
)

type app struct {
	lego.LogErr
}

func (a *app) GetName() string {
	return "Test App"
}

func (a *app) SetLogErr(logger lego.LogErr) {
	a.LogErr = logger
}

func (a *app) Register(p lego.Process) (io.Closer, error) {
	a.Info("application registered")

	return nil, nil
}

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	Run(ctx, &app{})
}
