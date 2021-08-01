package multilog

import (
	"github.com/vseinstrumentiru/lego/v2/log"
)

// Deprecated: use log.NewEntry
var NewNotification = log.NewEntry

// Deprecated: use log.NewErrEntry
var NewErrNotification = log.NewErrEntry

// Deprecated: use log.WrapOption
type WrapOption = log.WrapOption

// Deprecated: use log.WithFields
var WithFields = log.WithFields

// Deprecated: use log.WithDetails
var WithDetails = log.WithDetails

// Deprecated: use log.WrapEntry
var WrapNotification = log.WrapEntry
