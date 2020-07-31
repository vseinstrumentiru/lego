package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testDefaults struct{}

func (testDefaults) SetDefaults(e Env) {}

type testValidates struct{}

func (testValidates) Validate() error { return nil }

type testDefaultsValidates struct{}

func (testDefaultsValidates) SetDefaults(e Env) {}
func (testDefaultsValidates) Validate() error   { return nil }

func Test_ParseConfig(t *testing.T) {
	var test = struct {
		A1 testDefaults
		A2 testValidates
		A3 testDefaultsValidates
		A4 string
		A5 *struct {
			testDefaults
			testValidates
			B1 testDefaults
			B2 struct {
				C1 string
				testValidates
				C2 testValidates
			}
		}
	}{}

	defaults, validates := parseConfig(test, "app")
	assert.NotEmpty(t, defaults)
	assert.NotNil(t, defaults["app.A1"])
	assert.NotNil(t, defaults["app.A3"])
	assert.NotNil(t, defaults["app.A5"])
	assert.NotNil(t, defaults["app.A5.B1"])
	assert.NotEmpty(t, validates)
	assert.NotNil(t, validates["app.A2"])
	assert.NotNil(t, validates["app.A3"])
	assert.NotNil(t, validates["app.A5"])
	assert.NotNil(t, validates["app.A5.B2"])
	assert.NotNil(t, validates["app.A5.B2.C2"])
}
