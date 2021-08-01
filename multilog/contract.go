package multilog

import (
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/log"
)

// Deprecated: use logur.Level
type Level = logur.Level

// Deprecated: use log.Entry
type Entry = log.Entry

// Deprecated: use log.Logger
type Logger = log.Logger

// Deprecated: use log.EntryHandler
type EntryHandler = log.EntryHandler

// Deprecated: use log.ContextExtractor
type ContextExtractor = log.ContextExtractor

type (
	// Deprecated: use log.EntryMatcher
	EntryMatcher = log.EntryMatcher
	// Deprecated: use log.EntryErrMatcher
	EntryErrMatcher = log.EntryErrMatcher
)
