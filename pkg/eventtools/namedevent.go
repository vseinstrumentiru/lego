package eventtools

type NamedEvent interface {
	EventName() string
}

func EventName(v interface{}) string {
	if e, ok := v.(NamedEvent); ok {
		return e.EventName()
	}

	return ""
}
