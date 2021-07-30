package multilogprovider

import (
	"testing"

	zerolog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/multilog/log"
)

func TestProvide(t *testing.T) {
	c := container.New()

	ass := assert.New(t)

	ass.NoError(c.Register(func() logur.Logger {
		return zerologadapter.New(zerolog.Logger)
	}))
	ass.NoError(c.Instance(config.UndefinedApplication()))
	ass.NoError(c.Instance(&multilog.Config{Level: logur.Info}))
	ass.NoError(c.Instance(&log.Config{Stop: false}))
	ass.NoError(c.Register(Provide))

	ass.NoError(c.Execute(func(n multilog.Logger) {
		n.Info("test")
	}))
}
