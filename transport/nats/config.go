package nats

import (
	"time"

	"emperror.dev/errors"

	"github.com/vseinstrumentiru/lego/v2/config"
)

type Config struct {
	// Addr represents a single NATS server url to which the client
	// will be connecting. If the Servers option is also set, it
	// then becomes the first server in the Servers array.
	Addr string

	// Name is an optional name label which will be sent to the server
	// on CONNECT to identify the client.
	Name string

	// AllowReconnect enables reconnection logic to be used when we
	// encounter a disconnect from the current server.
	AllowReconnect bool

	// MaxReconnect sets the number of reconnect attempts that will be
	// tried before giving up. If negative, then it will never give up
	// trying to reconnect.
	MaxReconnect int

	// ReconnectWait sets the time to backoff after attempting a reconnect
	// to a server that we were already connected to previously.
	ReconnectWait time.Duration
}

func (c Config) SetDefaults(env config.Env) {
	env.SetDefault("allowReconnect", true)
	env.SetDefault("maxReconnect", -1)
	env.SetDefault("reconnectWait", "30s")
}

func (c Config) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("nats: addr is required"))
	}

	return
}
