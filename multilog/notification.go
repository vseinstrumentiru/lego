package multilog

import (
	"emperror.dev/errors"
	"logur.dev/logur"
)

func NewNotification(level Level, msg string, fields map[string]interface{}) Entry {
	return notification{
		msg:    msg,
		level:  level,
		fields: fields,
	}
}

type notification struct {
	msg     string
	level   Level
	fields  map[string]interface{}
	details []interface{}
}

func (n notification) Level() Level {
	return n.level
}

func (n notification) Message() string {
	return n.msg
}

func (n notification) Fields() map[string]interface{} {
	return n.fields
}

func (n notification) Details() []interface{} {
	return n.details
}

func (n notification) WithDetails(details ...interface{}) Entry {
	return WrapNotification(n, WithDetails(details...))
}

func (n notification) WithFields(fields map[string]interface{}) Entry {
	return WrapNotification(n, WithFields(fields))
}

func NewErrNotification(err error, fields map[string]interface{}) Entry {
	return errorNotification{
		err:    errors.WithStackDepthIf(err, 1),
		fields: fields,
	}
}

type errorNotification struct {
	err    error
	fields map[string]interface{}
}

func (n errorNotification) Level() Level {
	return logur.Error
}

func (n errorNotification) Error() string {
	return n.err.Error()
}

func (n errorNotification) Message() string {
	return n.err.Error()
}

func (n errorNotification) Fields() map[string]interface{} {
	return n.fields
}

func (n errorNotification) Details() []interface{} {
	return errors.GetDetails(n.err)
}

func (n errorNotification) WithDetails(details ...interface{}) Entry {
	n.err = errors.WithDetails(n.err, details...)

	return n
}

func (n errorNotification) Unwrap() error {
	return n.err
}

func (n errorNotification) WithFields(fields map[string]interface{}) Entry {
	for key, val := range fields {
		n.fields[key] = val
	}

	return n
}

type WrapOption func(n *wrapNotification)

func WithFields(fields map[string]interface{}) WrapOption {
	return func(n *wrapNotification) {
		for key, value := range fields {
			n.fields[key] = value
		}
	}
}

func WithDetails(details ...interface{}) WrapOption {
	return func(n *wrapNotification) {
		n.details = append(n.details, details...)
	}
}

func WrapNotification(parent Entry, opts ...WrapOption) Entry {
	entry := &wrapNotification{Entry: parent, fields: make(map[string]interface{})}

	for _, opt := range opts {
		opt(entry)
	}

	if len(entry.details) == 0 && len(entry.fields) == 0 {
		return parent
	}

	return entry
}

type wrapNotification struct {
	Entry
	details []interface{}
	fields  map[string]interface{}
}

func (n *wrapNotification) Fields() map[string]interface{} {
	fields := n.Entry.Fields()
	for key, val := range n.fields {
		fields[key] = val
	}

	return fields
}

func (n *wrapNotification) Details() []interface{} {
	details := n.Entry.Details()

	return append(details, n.details...)
}

func (n *wrapNotification) WithDetails(details ...interface{}) Entry {
	return WrapNotification(n, WithDetails(details...))
}

func (n *wrapNotification) WithFields(fields map[string]interface{}) Entry {
	return WrapNotification(n, WithFields(fields))
}
