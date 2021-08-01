module github.com/vseinstrumentiru/lego/v2

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/integrations/ocsql v0.1.7
	emperror.dev/emperror v0.33.0
	emperror.dev/errors v0.8.0
	github.com/AppsFlyer/go-sundheit v0.2.0
	github.com/Shopify/sarama v1.27.2
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-kafka v1.0.1
	github.com/ThreeDotsLabs/watermill-nats v1.0.5
	github.com/alecthomas/units v0.0.0-20201120081800-1786d5ef83d4 // indirect
	github.com/armon/go-metrics v0.3.6 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/cloudflare/tableflip v1.2.1
	github.com/dave/jennifer v1.4.1
	github.com/fatih/color v1.10.0 // indirect
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/getsentry/raven-go v0.2.0
	github.com/go-kit/kit v0.10.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.4
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/iancoleman/strcase v0.1.2
	github.com/jackc/pgx/v4 v4.10.1
	github.com/klauspost/compress v1.11.6 // indirect
	github.com/lebrains/gomongowrapper v0.0.0-20201208100026-f83272f79c09
	github.com/lithammer/shortuuid/v3 v3.0.5 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.4.0
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-streaming-server v0.20.0 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/stan.go v0.8.1
	github.com/newrelic/go-agent/v3 v3.9.0
	github.com/newrelic/go-agent/v3/integrations/nrgorilla v1.1.0
	github.com/newrelic/go-agent/v3/integrations/nrpkgerrors v1.0.0
	github.com/newrelic/newrelic-opencensus-exporter-go v0.4.0
	github.com/newrelic/newrelic-telemetry-sdk-go v0.5.1 // indirect
	github.com/nxadm/tail v1.4.6 // indirect
	github.com/oklog/run v1.1.0
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/statsd_exporter v0.18.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/sagikazarmark/appkit v0.10.0
	github.com/sagikazarmark/ocmux v0.2.0
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.5
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.0.0-20210108195828-e2f9c7f1fc8e
	google.golang.org/api v0.36.0 // indirect
	google.golang.org/genproto v0.0.0-20210108203827-ffc7fda8c3d7 // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gotest.tools v2.2.0+incompatible
	logur.dev/adapter/zerolog v0.5.0
	logur.dev/integration/watermill v0.5.0
	logur.dev/logur v0.17.0
	sagikazarmark.dev/mga v0.4.2
	sigs.k8s.io/controller-tools v0.4.1
)
