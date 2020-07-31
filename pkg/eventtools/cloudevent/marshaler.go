package cloudevent

import (
	"github.com/vseinstrumentiru/lego/tools/eventtools/cloudevent"
)

const (
	// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
	MetaCreatedAt = cloudevent.MetaCreatedAt
	// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
	MetaName = cloudevent.MetaName
	// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
	MetaSource = cloudevent.MetaSource
	// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
	MetaDataContentType = cloudevent.MetaDataContentType
)

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
type NamedEvent interface {
	cloudevent.NamedEvent
}

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
func EventName(v interface{}) string {
	return cloudevent.EventName(v)
}

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools/cloudevent
type Marshaller = cloudevent.Marshaller
