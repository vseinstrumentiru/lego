package lerr

const errNotFound = "not found"

// deprecated
func NotFound(details ...interface{}) error {
	return newErr(errNotFound, details)
}

// deprecated
func NotFoundWrap(err error, details ...interface{}) error {
	return wrap(err, errNotFound, details)
}
