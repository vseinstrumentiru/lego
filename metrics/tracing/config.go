package tracing

import (
	"strconv"
	"strings"

	"emperror.dev/errors"
	"go.opencensus.io/trace"
)

type Sampler func(trace.SamplingParameters) trace.SamplingDecision

const (
	always      = "always"
	probability = "probability"
	never       = "never"
)

type Config struct {
	Sampler Sampler

	// MaxAnnotationEventsPerSpan is max number of annotation events per span.
	MaxAnnotationEventsPerSpan int
	// MaxMessageEventsPerSpan is max number of message events per span.
	MaxMessageEventsPerSpan int
	// MaxAnnotationEventsPerSpan is max number of attributes per span.
	MaxAttributesPerSpan int
	// MaxLinksPerSpan is max number of links per span.
	MaxLinksPerSpan int
}

func (s *Sampler) UnmarshalText(b []byte) error {
	return s.FromString(b)
}

func (s *Sampler) Unmarshal(b []byte) error {
	return s.FromString(b)
}

func (s *Sampler) FromString(b []byte) error {
	text := string(b)

	if strings.HasPrefix(text, probability) {
		fractionTxt := text[len(probability)+1:]

		if fractionTxt == "" {
			return errors.NewWithDetails("probability trace without fraction", "use `sampler: probability:0.05`")
		}

		if fraction, err := strconv.ParseFloat(fractionTxt, 64); err != nil {
			return errors.WithStack(err)
		} else {
			*s = Sampler(trace.ProbabilitySampler(fraction))
			return nil
		}
	}

	switch text {
	case always:
		*s = Sampler(trace.AlwaysSample())
	case never:
		*s = Sampler(trace.NeverSample())
	default:
		return errors.NewWithDetails("undefined sampler type", "type", text)
	}

	return nil
}
