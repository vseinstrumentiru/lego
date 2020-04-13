package lerr

import (
	"emperror.dev/errors/match"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NotFoundMatch(t *testing.T) {
	err := NotFound()

	is := match.Is(NotFound()).MatchError(err)

	assert.True(t, is)
}

func Test_NotFoundWrapMatch(t *testing.T) {
	err := NotFound()

	is := match.Is(NotFound()).MatchError(err)

	assert.True(t, is)
}
