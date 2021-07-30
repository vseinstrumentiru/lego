package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

//nolint:gochecknoglobals
var (
	PublisherPublishTime = stats.Float64(
		"publish_time_seconds",
		"The time that a publishing attempt (success or not) took in seconds",
		stats.UnitMilliseconds,
	)
)

//nolint:gochecknoglobals
var (
	SubscriberReceivedMessage = stats.Float64(
		"subscriber_messages_received_total",
		"The total number of messages received by the subscriber",
		stats.UnitDimensionless,
	)
)

//nolint:gochecknoglobals
var (
	HandlerExecutionTime = stats.Float64(
		"handler_execution_time_seconds",
		"The total time elapsed while executing the handler function in seconds",
		stats.UnitMilliseconds,
	)
)

//nolint:gochecknoglobals
// The following tags are applied to stats recorded by this package.
var (
	HandlerName, _ = tag.NewKey("handler_name")

	PublisherName, _ = tag.NewKey("publisher_name")

	SubscriberName, _ = tag.NewKey("subscriber_name")

	Success, _ = tag.NewKey("success")

	Acked, _ = tag.NewKey("acked")
)

//nolint:gochecknoglobals
var (
	DefaultMillisecondsDistribution         = view.Distribution(0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
	DefaultMessageCountDistribution         = view.Distribution(1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536)
	DefaultHandlerExecutionTimeDistribution = view.Distribution(0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1)
)

//nolint:gochecknoglobals
var (
	PublisherPublishTimeView = &view.View{
		Name:        "publish_time_milliseconds",
		Description: "The time that a publishing attempt (success or not) took in milliseconds",
		Measure:     PublisherPublishTime,
		TagKeys:     []tag.Key{PublisherName, HandlerName, Success},
		Aggregation: DefaultMillisecondsDistribution,
	}

	SubscriberReceivedMessageView = &view.View{
		Name:        "subscriber_messages_received_total",
		Description: "The total number of messages received by the subscriber",
		Measure:     SubscriberReceivedMessage,
		TagKeys:     []tag.Key{SubscriberName, HandlerName, Acked},
		Aggregation: view.Count(),
	}

	HandlerExecutionTimeView = &view.View{
		Name:        "handler_execution_time_seconds",
		Description: "The total time elapsed while executing the handler function in seconds",
		Measure:     HandlerExecutionTime,
		TagKeys:     []tag.Key{HandlerName, Success},
		Aggregation: DefaultHandlerExecutionTimeDistribution,
	}
)
