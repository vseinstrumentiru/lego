package stan

import (
	"fmt"
	"os"
	"time"
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
	ClientID    string
	ClientIDGen ClientIDGen
	ClusterID   string
	// ConnectTimeout is the timeout for the initial Connect(). This value is also
	// used for some of the internal request/replies with the cluster.
	ConnectTimeout *time.Duration
	// AckTimeout is how long to wait when a message is published for an ACK from
	// the cluster. If the library does not receive an ACK after this timeout,
	// the Publish() call (or the AckHandler) will return ErrTimeout.
	AckTimeout *time.Duration
	// MaxPubAcksInflight specifies how many messages can be published without
	// getting ACKs back from the cluster before the Publish() or PublishAsync()
	// calls block.
	MaxPubAcksInflight *int
	// PingInterval is the interval at which client sends PINGs to the server
	// to detect the loss of a connection.
	PingInterval *int
	// PingMaxOut specifies the maximum number of PINGs without a corresponding
	// PONG before declaring the connection permanently lost.
	PingMaxOut *int
}

func (c *Config) GetClientID() string {
	if c.ClientIDGen == nil {
		return c.ClientID
	}

	return c.ClientIDGen(c.ClientID)
}
