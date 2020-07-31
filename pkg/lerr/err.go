package lerr

import "emperror.dev/errors"

// deprecated
type Typed interface {
	ErrorType() string
}

type leErr struct {
	error
	errType string
}

func newErr(errType string, msgs []interface{}) error {
	err := errors.NewPlain(errType)

	if len(msgs) > 0 {
		if msg, ok := msgs[0].(string); ok {
			err = errors.NewPlain(msg)
			msgs = msgs[1:]
		}

		err = errors.WithDetails(err, msgs...)
	}

	return leErr{error: errors.WithStackDepth(err, 2), errType: errType}
}

func wrap(err error, errType string, details []interface{}) error {
	if err == nil {
		return nil
	}

	return leErr{error: errors.WithDetails(errors.WithStackDepth(err, 2), details...), errType: errType}
}

func (e leErr) ErrorType() string {
	return e.errType
}

func (e leErr) Is(err error) bool {
	if lErr, ok := err.(Typed); ok {
		return lErr.ErrorType() == e.ErrorType()
	}

	return false
}
