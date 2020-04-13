package event

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"regexp"
	"time"
)

type Config struct {
	lego.WithSwitch `mapstructure:",squash"`

	Router   RouterConfig
	Provider string
	Nats     NatsProviderConfig
	Channel  GoChannelProviderConfig
}

const (
	NatsProvider    = "nats"
	ChannelProvider = "channel"
)

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.events.router.closeTimeout", time.Second)
	c.Nats.SetDefaults(env, flag)
}

func (c Config) Validate() (err error) {
	if !c.Enabled {
		return
	}

	if c.Provider == "" {
		err = errors.Append(err, errors.New("srv.events.provider is required"))
	}

	providerMap := map[string]struct{}{
		NatsProvider:    {},
		ChannelProvider: {},
	}

	if _, ok := providerMap[c.Provider]; !ok {
		err = errors.Append(err, errors.Errorf(`undefined srv.events.provider %s`, c.Provider))
	}

	if c.Provider == NatsProvider {
		err = errors.Append(err, c.Nats.Validate())
	}

	return
}

type RouterConfig struct {
	CloseTimeout time.Duration
}

type GoChannelProviderConfig struct {
	OutputChannelBuffer            int64
	Persistent                     bool
	BlockPublishUntilSubscriberAck bool
}

type NatsProviderConfig struct {
	Addr             string
	ClusterID        string
	ClientID         string
	QueueGroup       string
	DurableName      string
	SubscribersCount int
	CloseTimeout     time.Duration
	AckWaitTimeout   time.Duration
	CloudEvent       struct {
		lego.WithSwitch `mapstructure:",squash"`

		Source string
	} `mapstructure:"cloudevent"`
}

func (c NatsProviderConfig) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.events.nats.subscribersCount", 1)
	env.SetDefault("srv.events.nats.ackWaitTimeout", time.Second)
	env.SetDefault("srv.events.nats.closeTimeout", time.Second)
}

func (c NatsProviderConfig) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("srv.events.nats.addr is required"))
	}

	if c.ClusterID == "" {
		err = errors.Append(err, errors.New("srv.events.nats.clusterID is required"))
	}

	if c.ClientID == "" {
		err = errors.Append(err, errors.New("srv.events.nats.clientID is required"))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(c.ClientID) {
		err = errors.Append(err, errors.New("srv.events.nats.clientID should contain only alphanumeric characters, - or _"))
	}

	if c.CloudEvent.Enabled && c.CloudEvent.Source == "" {
		err = errors.Append(err, errors.New("srv.events.nats.cloutevent.source is required"))
	}

	return
}
