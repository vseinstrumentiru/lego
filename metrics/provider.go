package metrics

import (
	"fmt"
	"net/http"

	"emperror.dev/emperror"
	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"go.opencensus.io/zpages"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/common/netx"
	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type HealthArgs struct {
	dig.In
	Logger log.Logger
}

func ProvideHealthChecker(in HealthArgs) health.Health {
	logger := in.Logger.WithFields(map[string]interface{}{"component": "metrics.health"})

	healthz := health.New()
	healthz.WithCheckListener(NewLogCheckListener(logger))

	return healthz
}

type ServerArgs struct {
	dig.In
	Config *Config `optional:"true"`

	Health   health.Health
	Version  *version.Info
	Logger   log.Logger
	Pipeline *run.Group
	Upg      *tableflip.Upgrader `optional:"true"`
}

func ProvideMonitoringServer(in ServerArgs) *http.ServeMux {
	server := http.DefaultServeMux
	server.Handle("/version", version.Handler(in.Version))

	logger := in.Logger.WithFields(map[string]interface{}{"component": "metrics"})

	if in.Config == nil {
		in.Config = NewDefaultConfig()
	}

	if in.Config.Debug {
		zpages.Handle(server, "/debug")
	}

	{
		handler := healthhttp.HandleHealthJSON(in.Health)
		server.Handle("/healthz", handler)

		// Kubernetes style health checks
		server.HandleFunc("/healthz/live", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
		server.Handle("/healthz/ready", handler)
	}

	httpLn, err := netx.Listen("tcp", fmt.Sprintf(":%v", in.Config.Port), in.Upg)
	emperror.Panic(err)

	srv := &http.Server{
		Handler:  server,
		ErrorLog: log.NewErrorStandardLogger(logger),
	}

	in.Pipeline.Add(appkitrun.LogServe(logger.WithFields(map[string]interface{}{"port": in.Config.Port}))(
		appkitrun.HTTPServe(srv, httpLn, 0),
	))

	return server
}
