package LeGo

import (
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	"github.com/vseinstrumentiru/lego/internal/config"
	"github.com/vseinstrumentiru/lego/internal/monitor"
	"github.com/vseinstrumentiru/lego/internal/monitor/errorhandler"
	"github.com/vseinstrumentiru/lego/internal/monitor/log"
	"github.com/vseinstrumentiru/lego/pkg/build"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"net"
	"net/http"
	"os"
	"time"
)

type server struct {
	lego.LogErr
	Config    config.Config
	Telemetry *http.ServeMux
	Health    health.Health
	Upgrader  *tableflip.Upgrader
	Runner    run.Group
}

func newServer(name string, cfg lego.Config) *server {
	srv := &server{}
	var err error

	{
		srv.Config, err = config.Provide(cfg, config.WithDefaultName(name))
		logger := log.Provide(srv.Config.Monitor.Log, srv.Config.Env, srv.Config.Name)
		handler := errorhandler.Provide(srv.Config.Monitor.ErrorHandler, logger)

		srv.LogErr = lego.NewLogErr(logger, handler)

		srv.Info("starting application server", srv.Config.Build.Fields())

		if config.IsFileNotFound(err) {
			srv.Warn("configuration file not found")
		}

	}

	{
		err := srv.Config.Validate()
		if err != nil {
			srv.Error(err.Error())

			os.Exit(3)
		}
	}

	srv.Telemetry, srv.Health = monitor.Provide(srv, srv.Config.Monitor)

	srv.Upgrader, _ = tableflip.New(tableflip.Options{})

	return srv
}

func (s *server) Listen(network, addr string) (net.Listener, error) {
	return s.Upgrader.Listen(network, addr)
}

func (s *server) Background(execute func() error, interrupt func(error)) {
	s.Runner.Add(execute, interrupt)
}

func (s *server) ShutdownTimeout() time.Duration {
	return s.Config.ShutdownTimeout
}

func (s *server) Name() string {
	return s.Config.Name
}

func (s *server) Build() build.Info {
	return s.Config.Build
}

func (s *server) IsDebug() bool {
	return s.Config.Debug
}

func (s *server) DataCenterName() string {
	return s.Config.DataCenter
}
