package log

import (
	"logur.dev/logur"
)

func NewEntry(level logur.Level, msg string, fields map[string]interface{}) Entry {
	return &entry{
		msg:    msg,
		level:  level,
		fields: fields,
	}
}

type entry struct {
	msg     string
	level   logur.Level
	fields  map[string]interface{}
	details []interface{}
}

func (e *entry) Level() Level {
	return e.level
}

func (e *entry) Message() string {
	return e.msg
}

func (e *entry) Fields() map[string]interface{} {
	return e.fields
}

func (e *entry) Details() []interface{} {
	return e.details
}

func (e *entry) WithDetails(details ...interface{}) Entry {
	return WrapEntry(e, WithDetails(details...))
}

func (e *entry) WithFields(fields map[string]interface{}) Entry {
	return WrapEntry(e, WithFields(fields))
}
