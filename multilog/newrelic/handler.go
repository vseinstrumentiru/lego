package newrelic

import "github.com/vseinstrumentiru/lego/v2/log/handlers/newrelic"

// Deprecated: use newrelic.NewHandler
var Handler = newrelic.NewHandler
