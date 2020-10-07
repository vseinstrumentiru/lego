package events

import (
	"time"
)

type Config struct {
	AckTimeout       *time.Duration
	ReconnectTimeout *time.Duration
	CloseTimeout     *time.Duration
}
