package logger_test

import (
	"bytes"
	"testing"

	zerolog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/log/handlers/console"
	"github.com/vseinstrumentiru/lego/v2/log/logger"
)

func newTestLoggerWithConsole(level logur.Level) (log.Logger, *bytes.Buffer) {
	c := container.New()
	out := &bytes.Buffer{}
	c.Instance(&logger.Config{Level: level})
	c.Instance(console.CustomWriterConfig(out))
	c.Register(console.Provide)
	c.Register(logger.Provide)

	var l log.Logger
	c.Execute(func(n log.Logger) {
		l = n
	})

	return l, out
}

func TestProvide(t *testing.T) {
	c := container.New()

	ass := assert.New(t)

	ass.NoError(c.Instance(config.UndefinedApplication()))
	ass.NoError(c.Instance(&logger.Config{Level: logur.Info}))
	out := &bytes.Buffer{}
	ass.NoError(c.Instance(console.CustomWriterConfig(out)))
	ass.NoError(c.Register(console.Provide))
	ass.NoError(c.Register(logger.Provide))

	ass.NoError(c.Execute(func(n log.Logger) {
		n.Info("test")
	}))
}

func TestCustomConsoleLogger(t *testing.T) {
	c := container.New()

	ass := assert.New(t)

	ass.NoError(c.Register(func() logur.Logger {
		return zerologadapter.New(zerolog.Logger)
	}))
	ass.NoError(c.Instance(&logger.Config{Level: logur.Trace}))
	out := &bytes.Buffer{}
	ass.NoError(c.Instance(console.CustomWriterConfig(out)))
	ass.NoError(c.Register(console.Provide))
	ass.NoError(c.Register(logger.Provide))

	ass.NoError(c.Execute(func(n log.Logger) {
		n.Info("test")
	}))
}

func TestConsole(t *testing.T) {
	l, out := newTestLoggerWithConsole(logur.Trace)
	l.Info("Some long message exist and it's ok")
	assert.Contains(t, out.String(), "message exist")
}

func TestLevel(t *testing.T) {
	l, out := newTestLoggerWithConsole(logur.Trace)

	el := l.WithLevel(logur.Error)

	el.Info("Some long message exist and it's ok")
	assert.NotContains(t, out.String(), "message exist")

	el.Error("Some long message exist and it's ok")
	assert.Contains(t, out.String(), "message exist")
}
