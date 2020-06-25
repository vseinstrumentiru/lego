package lerr

import "strings"

type likeErr struct {
	msg string
}

func NewLike(msg string) error {
	return likeErr{msg: msg}
}

func (e likeErr) Error() string {
	panic("this error is used only for check another errors")
}

func (e likeErr) Is(err error) bool {
	return Like(err, e.msg)
}

func Like(err error, msg string) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), msg)
}
