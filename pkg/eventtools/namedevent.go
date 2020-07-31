package eventtools

import "github.com/vseinstrumentiru/lego/tools/eventtools"

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools
type NamedEvent interface {
	eventtools.NamedEvent
}

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/eventtools
func EventName(v interface{}) string {
	return eventtools.EventName(v)
}
