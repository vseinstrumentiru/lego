package deprecated

import (
	"fmt"
	"time"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters"
	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/multilog/log"
	"github.com/vseinstrumentiru/lego/v2/multilog/sentry"
	"github.com/vseinstrumentiru/lego/v2/transport/grpc"
	"github.com/vseinstrumentiru/lego/v2/transport/http"
	"github.com/vseinstrumentiru/lego/v2/transport/nats"
	"github.com/vseinstrumentiru/lego/v2/transport/stan"
)

// Deprecated: use LeGo V2
type Config interface {
	config.Validatable
	config.WithDefaults
}

type FullConfig struct {
	App interface{}
	Srv struct {
		AppConfig config.Application `mapstructure:",squash" load:"true"`
		Debug     bool
		Http      *http.Config
		Grpc      *grpc.Config
		Events    *struct {
			DefaultProvider string
			Providers       struct {
				Nats map[string]struct {
					Addr              string
					ClusterID         string
					ClientID          string
					ClientIDSuffixGen string
					QueueGroup        string
					DurableName       string
					CloseTimeout      time.Duration
					AckWaitTimeout    time.Duration
				}
				Channel map[string]gochannel.Config
			}
		}
		Monitor *struct {
			Log          log.Config
			Errorhandler *struct {
				Sentry struct {
					DSN string
				}
			}
			Exporter *struct {
				Jaeger *struct {
					Addr string
				}

				NewRelic *struct {
					Key string
				}
			}
			Trace *struct {
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
		}
		ShutdownTimeout time.Duration
	}
}

func (c *FullConfig) Convert() []interface{} {
	srv := c.Srv
	var configs []interface{}

	if srv.Events != nil {
		if len(srv.Events.Providers.Nats) != 0 {
			for _, i := range srv.Events.Providers.Nats {
				if i.Addr == "" {
					continue
				}
				{
					configs = append(configs, &nats.Config{
						Addr:           i.Addr,
						AllowReconnect: true,
						MaxReconnect:   -1,
						ReconnectWait:  30 * time.Second,
					})
				}
				{
					clientIdGen := stan.NoGeneration
					if i.ClientIDSuffixGen != "" {
						clientIdGen = stan.HostSuffixed
					}
					configs = append(configs, &stan.Config{
						ClientID:           i.ClientID,
						ClusterID:          i.ClientID,
						GroupName:          i.QueueGroup,
						DurableName:        i.DurableName,
						ClientIDGen:        clientIdGen,
						ConnectTimeout:     30 * time.Second,
						MaxPubAcksInflight: 16384,
						PingInterval:       10,
						PingMaxOut:         20,
					})
				}

				break
			}
		}

		if len(srv.Events.Providers.Channel) != 0 {
			for _, i := range srv.Events.Providers.Channel {
				configs = append(configs, i)
				break
			}
		}
	}

	if srv.Monitor != nil {
		if srv.Monitor.Errorhandler != nil {
			if srv.Monitor.Errorhandler.Sentry.DSN != "" {
				configs = append(configs, &sentry.Config{
					Addr: srv.Monitor.Errorhandler.Sentry.DSN,
				})
			}
		}

		if srv.Monitor.Exporter != nil {
			if srv.Monitor.Exporter.Jaeger != nil && srv.Monitor.Exporter.Jaeger.Addr != "" {
				configs = append(configs, &exporters.Jaeger{Addr: srv.Monitor.Exporter.Jaeger.Addr})
			}

			if srv.Monitor.Exporter.NewRelic != nil && srv.Monitor.Exporter.NewRelic.Key != "" {
				configs = append(configs, &exporters.NewRelic{
					Enabled:          true,
					TelemetryEnabled: true,
					Key:              srv.Monitor.Exporter.NewRelic.Key,
				})
			}
		}

		if srv.Monitor.Trace != nil && srv.Monitor.Trace.Sampling.Sampler != "" {
			var sampler tracing.Sampler
			if srv.Monitor.Trace.Sampling.Sampler == "probability" {
				emperror.Panic(sampler.FromString([]byte(srv.Monitor.Trace.Sampling.Sampler + ":" + fmt.Sprint(srv.Monitor.Trace.Sampling.Fraction))))
			} else {
				emperror.Panic(sampler.FromString([]byte(srv.Monitor.Trace.Sampling.Sampler)))
			}

			configs = append(configs, &tracing.Config{
				Sampler:                    sampler,
				MaxAnnotationEventsPerSpan: srv.Monitor.Trace.MaxAnnotationEventsPerSpan,
				MaxMessageEventsPerSpan:    srv.Monitor.Trace.MaxMessageEventsPerSpan,
				MaxAttributesPerSpan:       srv.Monitor.Trace.MaxAttributesPerSpan,
				MaxLinksPerSpan:            srv.Monitor.Trace.MaxLinksPerSpan,
			})
		}
	}

	return configs
}
