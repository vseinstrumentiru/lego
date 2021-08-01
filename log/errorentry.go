package log

import (
	"emperror.dev/errors"
	"logur.dev/logur"
)

func NewErrEntry(err error, fields map[string]interface{}) Entry {
	return &errorEntry{
		err:    errors.WithStackDepthIf(err, 1),
		fields: fields,
	}
}

type errorEntry struct {
	err    error
	fields map[string]interface{}
}

func (e *errorEntry) Level() logur.Level {
	return logur.Error
}

func (e *errorEntry) Error() string {
	return e.err.Error()
}

func (e *errorEntry) Message() string {
	return e.err.Error()
}

func (e *errorEntry) Fields() map[string]interface{} {
	return e.fields
}

func (e *errorEntry) Details() []interface{} {
	return errors.GetDetails(e.err)
}

func (e *errorEntry) WithDetails(details ...interface{}) Entry {
	e.err = errors.WithDetails(e.err, details...)

	return e
}

func (e *errorEntry) Unwrap() error {
	return e.err
}

func (e *errorEntry) WithFields(fields map[string]interface{}) Entry {
	for key, val := range fields {
		e.fields[key] = val
	}

	return e
}
