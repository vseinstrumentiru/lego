package event

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"regexp"
	"time"
)

type Config struct {
	lego2.WithSwitch `mapstructure:",squash"`
	RouterConfig     `mapstructure:",squash"`

	DefaultProvider string

	Router struct {
		Pub map[string]string
		Sub map[string]string
	}

	Providers struct {
		Nats    map[string]NatsProviderConfig
		Channel map[string]GoChannelProviderConfig
	}
}

type Provider struct {
	Pub bool
	Sub bool
}

const (
	NatsProvider    = "nats"
	ChannelProvider = "channel"
)

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.events.router.closeTimeout", time.Second)
	env.SetDefault("srv.events.router.maxRetries", 1)
	env.SetDefault("srv.events.router.maxRetryInterval", time.Millisecond*200)

	for name, cfg := range c.Providers.Nats {
		cfg.SetDefaults("srv.events.provider.nats."+name, env, flag)
	}

	for name, cfg := range c.Providers.Channel {
		cfg.SetDefaults("srv.events.provider.channel."+name, env, flag)
	}
}

func (c Config) Validate() (err error) {
	if !c.Enabled {
		return
	}

	if c.DefaultProvider == "" {
		err = errors.Append(err, errors.New("srv.events.defaultProvider is required"))
	}

	providerMap := make(map[string]Provider)

	for name, cfg := range c.Providers.Nats {
		if _, ok := providerMap[name]; ok {
			err = errors.Append(err, errors.Errorf(`provider name %s already exist`, name))
		} else {
			providerMap[name] = cfg.Provider
		}
		err = errors.Append(err, cfg.Validate())
	}

	for name, cfg := range c.Providers.Channel {
		if _, ok := providerMap[name]; ok {
			err = errors.Append(err, errors.Errorf(`provider name %s already exist`, name))
		} else {
			providerMap[name] = Provider{
				Pub: true,
				Sub: true,
			}
		}
		err = errors.Append(err, cfg.Validate())
	}

	if _, ok := providerMap[c.DefaultProvider]; !ok {
		err = errors.Append(err, errors.Errorf(`default provider with name %s not found`, c.DefaultProvider))
	}

	for topic, provider := range c.Router.Sub {
		if p, ok := providerMap[provider]; !ok {
			err = errors.Append(err, errors.Errorf(`provider for topic %s with name %s not found`, topic, provider))
		} else if !p.Sub {
			err = errors.Append(err, errors.Errorf(`provider for topic %s with name %s has not subscriber`, topic, provider))
		}
	}

	for topic, provider := range c.Router.Pub {
		if p, ok := providerMap[provider]; !ok {
			err = errors.Append(err, errors.Errorf(`provider for topic %s with name %s not found`, topic, provider))
		} else if !p.Pub {
			err = errors.Append(err, errors.Errorf(`provider for topic %s with name %s has not publisher`, topic, provider))
		}
	}

	return
}

type RouterConfig struct {
	CloseTimeout     time.Duration
	MaxRetries       int
	MaxRetryInterval time.Duration
}

type GoChannelProviderConfig struct {
	Provider                       `mapstructure:",squash"`
	OutputChannelBuffer            int64
	Persistent                     bool
	BlockPublishUntilSubscriberAck bool
}

func (g GoChannelProviderConfig) Validate() error {
	return nil
}

func (g GoChannelProviderConfig) SetDefaults(key string, env *viper.Viper, flag *pflag.FlagSet) {
	return
}

const (
	natsClientIDSuffixHost = "host"
)

type NatsProviderConfig struct {
	Provider          `mapstructure:",squash"`
	Addr              string
	ClusterID         string
	ClientID          string
	ClientIDSuffixGen string
	QueueGroup        string
	DurableName       string
	SubscribersCount  int
	CloseTimeout      time.Duration
	AckWaitTimeout    time.Duration
	PanicOnLost       bool
	CloudEvent        struct {
		lego2.WithSwitch `mapstructure:",squash"`

		Source string
	} `mapstructure:"cloudevent"`
}

func (c NatsProviderConfig) SetDefaults(path string, env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault(path+".pub", true)
	env.SetDefault(path+".sub", false)
	env.SetDefault(path+".subscribersCount", 1)
	env.SetDefault(path+".ackWaitTimeout", time.Second)
	env.SetDefault(path+".closeTimeout", time.Second)
	env.SetDefault(path+".clientIDSuffixGen", natsClientIDSuffixHost)
}

func (c NatsProviderConfig) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("addr is required in nats provider"))
	}

	if c.ClusterID == "" {
		err = errors.Append(err, errors.New("clusterID is required in nats provider"))
	}

	if c.ClientID == "" {
		err = errors.Append(err, errors.New("clientID is required in nats provider"))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(c.ClientID) {
		err = errors.Append(err, errors.New("clientID should contain only alphanumeric characters, - or _  in nats provider"))
	}

	if c.CloudEvent.Enabled && c.CloudEvent.Source == "" {
		err = errors.Append(err, errors.New("cloutevent.source is required in nats provider"))
	}

	return
}
