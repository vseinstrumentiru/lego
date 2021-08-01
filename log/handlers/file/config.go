package file

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"logur.dev/logur"
)

type Config struct {
	Level             logur.Level
	Stop              bool
	lumberjack.Logger `mapstructure:",squash"`
}
