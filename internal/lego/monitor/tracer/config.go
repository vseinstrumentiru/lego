package tracer

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	"go.opencensus.io/trace"
	"strings"
)

type Config struct {
	lego2.WithSwitch `mapstructure:",squash"`

	// Sampling describes the default sampler used when creating new spans.
	Sampling struct {
		Sampler  string
		Fraction float64
	}
	// MaxAnnotationEventsPerSpan is max number of annotation events per span.
	MaxAnnotationEventsPerSpan int
	// MaxMessageEventsPerSpan is max number of message events per span.
	MaxMessageEventsPerSpan int
	// MaxAnnotationEventsPerSpan is max number of attributes per span.
	MaxAttributesPerSpan int
	// MaxLinksPerSpan is max number of links per span.
	MaxLinksPerSpan int
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.monitor.trace.sampling.sampler", "never")
}

func (c Config) Validate() (err error) {
	if !c.Enabled {
		return
	}

	if c.Sampling.Sampler != "always" && c.Sampling.Sampler != "never" && c.Sampling.Sampler != "probability" {
		err = errors.Append(err, errors.New("srv.monitor.trace.sampling.sampler must be on of [always, never, probability]"))
	}

	return
}

func (c Config) Config() trace.Config {
	config := trace.Config{
		MaxAnnotationEventsPerSpan: c.MaxAnnotationEventsPerSpan,
		MaxMessageEventsPerSpan:    c.MaxMessageEventsPerSpan,
		MaxAttributesPerSpan:       c.MaxAttributesPerSpan,
		MaxLinksPerSpan:            c.MaxLinksPerSpan,
	}

	switch strings.ToLower(strings.TrimSpace(c.Sampling.Sampler)) {
	case "always":
		config.DefaultSampler = trace.AlwaysSample()

	case "never":
		config.DefaultSampler = trace.NeverSample()

	case "probability":
		config.DefaultSampler = trace.ProbabilitySampler(c.Sampling.Fraction)
	}

	return config
}
