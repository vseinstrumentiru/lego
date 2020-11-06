package tracing

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"

	"github.com/vseinstrumentiru/lego/internal/env"
)

func TestSampler_ProbabilityUnmarshall(t *testing.T) {
	var c Config
	err := json.Unmarshal([]byte(`{"Sampler":"probability:1"}`), &c)
	ass := assert.New(t)
	ass.NoError(err)
	_, span := trace.StartSpan(context.Background(), "test")
	defer span.End()
	params := trace.SamplingParameters{
		ParentContext:   span.SpanContext(),
		TraceID:         trace.TraceID{},
		SpanID:          trace.SpanID{},
		Name:            "ttt",
		HasRemoteParent: false,
	}
	ass.Equal(Sampler(trace.ProbabilitySampler(1))(params), c.Sampler(params))
}

func TestSampler_AlwaysUnmarshall(t *testing.T) {
	var c Config
	err := json.Unmarshal([]byte(`{"Sampler":"always"}`), &c)
	ass := assert.New(t)
	ass.NoError(err)
	_, span := trace.StartSpan(context.Background(), "test")
	defer span.End()
	params := trace.SamplingParameters{
		ParentContext:   span.SpanContext(),
		TraceID:         trace.TraceID{},
		SpanID:          trace.SpanID{},
		Name:            "ttt",
		HasRemoteParent: false,
	}
	ass.Equal(Sampler(trace.AlwaysSample())(params), c.Sampler(params))
}

func TestSampler_NeverUnmarshall(t *testing.T) {
	var c Config
	err := json.Unmarshal([]byte(`{"Sampler":"never"}`), &c)
	ass := assert.New(t)
	ass.NoError(err)
	_, span := trace.StartSpan(context.Background(), "test")
	defer span.End()
	params := trace.SamplingParameters{
		ParentContext:   span.SpanContext(),
		TraceID:         trace.TraceID{},
		SpanID:          trace.SpanID{},
		Name:            "ttt",
		HasRemoteParent: false,
	}
	ass.Equal(Sampler(trace.NeverSample())(params), c.Sampler(params))
}

func TestSampler_WrongUnmarshall(t *testing.T) {
	var c Config
	err := json.Unmarshal([]byte(`{"Sampler":"wrong"}`), &c)
	ass := assert.New(t)
	ass.Error(err)
}

func TestSampler_ViperUnmarshall(t *testing.T) {
	_ = os.Setenv("CFG_SAMPLER", "probability:0.5")
	var c Config

	e := env.NewConfigEnv("cfg").(interface{ Configure(cfg env.Config) error })
	err := e.Configure(&c)

	ass := assert.New(t)
	ass.NoError(err)
	ass.NotNil(c.Sampler)
}
