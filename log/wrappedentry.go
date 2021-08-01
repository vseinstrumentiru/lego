package log

type WrapOption func(n *wrappedEntry)

func WithFields(fields map[string]interface{}) WrapOption {
	return func(n *wrappedEntry) {
		for key, value := range fields {
			n.fields[key] = value
		}
	}
}

func WithDetails(details ...interface{}) WrapOption {
	return func(n *wrappedEntry) {
		n.details = append(n.details, details...)
	}
}

func WrapEntry(parent Entry, opts ...WrapOption) Entry {
	e := &wrappedEntry{Entry: parent, fields: make(map[string]interface{})}

	for _, opt := range opts {
		opt(e)
	}

	if len(e.details) == 0 && len(e.fields) == 0 {
		return parent
	}

	return e
}

type wrappedEntry struct {
	Entry
	details []interface{}
	fields  map[string]interface{}
}

func (e *wrappedEntry) Fields() map[string]interface{} {
	fields := e.Entry.Fields()
	for key, val := range e.fields {
		fields[key] = val
	}

	return fields
}

func (e *wrappedEntry) Details() []interface{} {
	details := e.Entry.Details()

	return append(details, e.details...)
}

func (e *wrappedEntry) WithDetails(details ...interface{}) Entry {
	return WrapEntry(e, WithDetails(details...))
}

func (e *wrappedEntry) WithFields(fields map[string]interface{}) Entry {
	return WrapEntry(e, WithFields(fields))
}
