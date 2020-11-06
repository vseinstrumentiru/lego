package events

import "time"

type RouterConfig struct {
	// CloseTimeout determines how long router should work for handlers when closing.
	CloseTimeout time.Duration
	// Retry provides a middleware that retries the handler if errors are returned.
	// The retry behaviour is configurable, with exponential backoff and maximum elapsed time.
	// if retries limit was exceeded, message is sent to poison queue (poison_queue topic)
	Retry *struct {
		Count    int
		Interval time.Duration
	}
	// CorrelationID adds correlation ID to all messages produced by the handler.
	// ID is based on ID from message received by handler.
	Correlation bool

	// Throttle provides a middleware that limits the amount of messages processed per unit of time.
	// This may be done e.g. to prevent excessive load caused by running a handler on a long queue of unprocessed messages.
	Throttle *struct {
		Count    int64
		Interval time.Duration
	}
	// Timeout makes the handler cancel the incoming message's context after a specified time.
	// Any timeout-sensitive functionality of the handler should listen on msg.Context().Done() to know when to fail.
	Timeout *struct {
		Interval time.Duration
	}
}
