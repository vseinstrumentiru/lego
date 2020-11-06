package sentry

import "logur.dev/logur"

type Config struct {
	Addr  string
	Level logur.Level
	Stop  bool
}
