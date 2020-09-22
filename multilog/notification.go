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
	if len(details) == 0 {
		return n
	}

	return notification{
		msg:     n.msg,
		level:   n.level,
		fields:  n.fields,
		details: append(n.details, details...),
	}
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
