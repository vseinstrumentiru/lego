package cloudevent

import "time"

type EventWithTime interface {
	SetCreatedAt(t time.Time)
	GetCreatedAt() time.Time
}
