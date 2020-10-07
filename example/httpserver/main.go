package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/multilog"
	"github.com/vseinstrumentiru/lego/multilog/log"
	"github.com/vseinstrumentiru/lego/server"
	httpcfg "github.com/vseinstrumentiru/lego/transport/http"
)

type Config struct {
	config.Application `mapstructure:",squash"`
	Http               httpcfg.Config
	Logger             multilog.Config
	Log                log.Config
}

type Tmp struct {
	Val string
}

type App struct {
	b      string
	tmp    Tmp
	Logger multilog.Logger
}

func (a App) ConfigureRoutes(r *mux.Router) {
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
}

var app App = App{b: "test"}

func main() {
	server.Run(app, &Config{})
}
