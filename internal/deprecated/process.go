package deprecated

import (
	"net"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/multilog"
	"github.com/vseinstrumentiru/lego/version"
)

type Process interface {
	LogErr

	Name() string
	DataCenterName() string
	Build() version.Info
	Env() string
	IsDebug() bool

	Listen(network, addr string) (net.Listener, error)
	Background(execute func() error, interrupt func(error))
	RegisterCheck(cfg *health.Config) error
	ShutdownTimeout() time.Duration
}

type processArgs struct {
	dig.In
	App    *config.Application
	Ver    version.Info
	Checks health.Health
	Upg    *tableflip.Upgrader
	Run    *run.Group
	Log    multilog.Logger
}

func NewProcess(args processArgs) Process {
	return &process{
		LogErr: &multilogWrapper{Logger: args.Log},
		app:    *args.App,
		ver:    args.Ver,
		checks: args.Checks,
		upg:    args.Upg,
		run:    args.Run,
	}
}

type process struct {
	LogErr
	app    config.Application
	ver    version.Info
	checks health.Health
	upg    *tableflip.Upgrader
	run    *run.Group
}

func (p *process) Name() string {
	return p.app.Name
}

func (p *process) DataCenterName() string {
	return p.app.DataCenter
}

func (p *process) Build() version.Info {
	return p.ver
}

func (p *process) Env() string {
	return ""
}

func (p *process) IsDebug() bool {
	return false
}

func (p *process) Listen(network, addr string) (net.Listener, error) {
	return p.upg.Listen(network, addr)
}

func (p *process) Background(execute func() error, interrupt func(error)) {
	p.run.Add(execute, interrupt)
}

func (p *process) RegisterCheck(cfg *health.Config) error {
	return p.checks.RegisterCheck(cfg)
}

func (p *process) ShutdownTimeout() time.Duration {
	return 30 * time.Second
}
