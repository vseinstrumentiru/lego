package core

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_Component(t *testing.T) {
	type tmp struct {
		CoreComponent
	}

	name := componentName(tmp.Component)

	assert.Equal(t, "core", name)
}

func Test_Type_Wrong(t *testing.T) {
	type tmp struct {
	}

	name := typ(tmp{})

	assert.Equal(t, "", name)
}
