package lerr

const errNotFound = "not found"

func NotFound(details ...interface{}) error {
	return newErr(errNotFound, details)
}

func NotFoundWrap(err error, details ...interface{}) error {
	return wrap(err, errNotFound, details)
}
