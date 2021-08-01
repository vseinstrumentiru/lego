package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/vseinstrumentiru/lego/v2/log"
)

func NewLogger(l log.Logger) resty.Logger {
	return &logger{Logger: l}
}

type logger struct {
	log.Logger
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}
