package sentry

import "github.com/vseinstrumentiru/lego/v2/log/handlers/sentry"

// Deprecated: use sentry.NewHandler
var Handler = sentry.NewHandler
