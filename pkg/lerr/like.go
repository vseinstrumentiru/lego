package lerr

import (
	"strings"
)

func Like(err error, msg string) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), msg)
}
