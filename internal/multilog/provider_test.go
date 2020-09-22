package multilog

import (
	"testing"

	zerolog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/components/notify/log"

	container2 "github.com/vseinstrumentiru/lego/internal/container"
)

func TestProvide(t *testing.T) {
	c := container2.New()

	ass := assert.New(t)

	ass.NoError(c.Register(func() logur.Logger {
		return zerologadapter.New(zerolog.Logger)
	}))
	ass.NoError(c.Instance(&multilog.Config{Level: logur.Info}))
	ass.NoError(c.Instance(&log.Config{Stop: false}))
	ass.NoError(c.Register(Provide))

	ass.NoError(c.Execute(func(n multilog.Multilog) {
		n.Info("test")
	}))
}
