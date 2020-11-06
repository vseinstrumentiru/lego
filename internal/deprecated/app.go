package deprecated

import (
	"io"
	"reflect"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/vseinstrumentiru/lego/v2/config"
	di "github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
)

func NewApp(dApp App) (interface{}, *FullConfig) {
	app := &app{
		i: dApp,
	}

	if cApp, ok := dApp.(interface{ GetConfig() Config }); ok {
		appCfg := cApp.GetConfig()
		val := reflect.ValueOf(appCfg)
		if val.Kind() == reflect.Ptr {
			app.cfg.App = val.Interface()
		} else if val.CanAddr() {
			app.cfg.App = val.Addr().Interface()
		} else {
			val = reflect.New(val.Type())
			app.cfg.App = val.Interface()
		}
	}

	app.cfg.Srv.AppConfig.Name = dApp.GetName()

	return app, &app.cfg
}

type App interface {
	GetName() string
	SetLogErr(logErr LogErr)
}

type app struct {
	i   App
	cfg FullConfig
}

func (app app) Providers() []interface{} {
	return []interface{}{
		NewEventManager,
		NewProcess,
	}
}

func (app app) Configure(c di.Container, log multilog.Logger, close *shutdown.CloseGroup) error {
	configs := app.cfg.Convert()

	app.i.SetLogErr(&multilogWrapper{Logger: log.WithFields(map[string]interface{}{"component": "app"})})

	for i := 0; i < len(configs); i++ {
		if v, ok := configs[i].(config.Validatable); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}

		if err := c.Instance(configs[i]); err != nil {
			return err
		}
	}

	if appCfg, ok := app.i.(interface{ SetConfig(config Config) }); ok && app.cfg.App != nil {
		val := reflect.ValueOf(app.cfg.App).Elem()
		appCfg.SetConfig(val.Interface().(Config))
	}

	if err := c.Make(app.i); err != nil {
		return err
	}

	if appPub, ok := app.i.(interface {
		RegisterEventDispatcher(publisher message.Publisher) error
	}); ok {
		type args struct {
			dig.In
			Nats    *nats.StreamingPublisher `optional:"true"`
			Channel *gochannel.GoChannel     `optional:"true"`
		}
		err := c.Execute(func(in args) error {
			if in.Nats != nil {
				return appPub.RegisterEventDispatcher(in.Nats)
			} else if in.Channel != nil {
				return appPub.RegisterEventDispatcher(in.Channel)
			}
			return errors.New("publisher not found")
		})

		if err != nil {
			return err
		}
	}

	if appSub, ok := app.i.(interface{ RegisterEventHandlers(em EventManager) error }); ok {
		err := c.Execute(func(em EventManager) error {
			return appSub.RegisterEventHandlers(em)
		})

		if err != nil {
			return err
		}
	}

	if appProc, ok := app.i.(interface {
		Register(p Process) (io.Closer, error)
	}); ok {
		err := c.Execute(func(p Process) error {
			off, err := appProc.Register(p)

			if off != nil {
				close.Add(off)
			}

			return err
		})

		if err != nil {
			return err
		}
	}

	if grpcApp, ok := app.i.(interface {
		RegisterGRPC(server *grpc.Server) error
	}); ok {
		err := c.Execute(func(server *grpc.Server) error {
			return grpcApp.RegisterGRPC(server)
		})

		if err != nil {
			return err
		}
	}

	if httpApp, ok := app.i.(interface {
		RegisterHTTP(router *mux.Router) error
	}); ok {
		err := c.Execute(func(router *mux.Router) error {
			return httpApp.RegisterHTTP(router)
		})

		if err != nil {
			return err
		}
	}

	if runApp, ok := app.i.(interface {
		Run(terminate chan bool) error
	}); ok {
		err := c.Execute(func(run run.Group) {
			terminate := make(chan bool, 1)
			run.Add(func() error {
				return runApp.Run(terminate)
			}, func(err error) {
				terminate <- true
			})
		})

		if err != nil {
			return err
		}
	}

	return nil
}
