package lerr

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LikeError(t *testing.T) {
	err := bytes.ErrTooLarge

	assert.True(t, Is(NewLike("too large"), err))
}
