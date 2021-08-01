package newrelicexporter

import (
	"github.com/vseinstrumentiru/lego/v2/log"
)

type loggerWrap struct {
	log.Logger
}

func (l loggerWrap) Error(msg string, context map[string]interface{}) {
	l.Logger.Error(msg, context)
}

func (l loggerWrap) Warn(msg string, context map[string]interface{}) {
	l.Logger.Warn(msg, context)
}

func (l loggerWrap) Info(msg string, context map[string]interface{}) {
	l.Logger.Info(msg, context)
}

func (l loggerWrap) Debug(msg string, context map[string]interface{}) {
	l.Logger.Debug(msg, context)
}

func (l loggerWrap) DebugEnabled() bool {
	return false
}
