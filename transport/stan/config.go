package stan

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"emperror.dev/errors"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/events"
)

const (
	hostGen = "host"
)

type ClientIDGen func(clientID string) string

func NoGeneration(clientID string) string {
	return clientID
}

func HostSuffixed(clientID string) string {
	host, _ := os.Hostname()

	return fmt.Sprintf("%s_%s", clientID, host)
}

func (s *ClientIDGen) UnmarshalText(b []byte) error {
	switch string(b) {
	case hostGen:
		*s = HostSuffixed
	default:
		*s = NoGeneration
	}

	return nil
}

type Config struct {
	events.Config `mapstructure:",squash"`

	ClientID    string
	ClientIDGen ClientIDGen
	ClusterID   string
	GroupName   string
	DurableName string
	// ConnectTimeout is the timeout for the initial Connect(). This value is also
	// used for some of the internal request/replies with the cluster.
	ConnectTimeout time.Duration
	// MaxPubAcksInflight specifies how many messages can be published without
	// getting ACKs back from the cluster before the Publish() or PublishAsync()
	// calls block.
	MaxPubAcksInflight int
	// PingInterval is the interval at which client sends PINGs to the server
	// to detect the loss of a connection.
	PingInterval int
	// PingMaxOut specifies the maximum number of PINGs without a corresponding
	// PONG before declaring the connection permanently lost.
	PingMaxOut int
}

func (c *Config) SetDefaults(env config.Env) {
	env.SetDefault("ackTimeout", "30s")
	env.SetDefault("connectTimeout", "5s")
	env.SetDefault("maxPubAcksInflight", 16384)
	env.SetDefault("pingInterval", 10)
	env.SetDefault("pingMaxOut", 20)
}

func (c *Config) Validate() (err error) {
	if c.ClusterID == "" {
		err = errors.Append(err, errors.New("stan: clusterID is required"))
	}

	if c.ClientID == "" {
		err = errors.Append(err, errors.New("stan: clientID is required in"))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(c.ClientID) {
		err = errors.Append(err, errors.New("stan: clientID should contain only alphanumeric characters, - or _"))
	}

	return
}

func (c *Config) GetClientID() string {
	if c.ClientIDGen == nil {
		return c.ClientID
	}

	return c.ClientIDGen(c.ClientID)
}
