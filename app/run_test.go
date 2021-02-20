package app

import (
	"testing"
)

type testApp struct {
}

type testConfig struct {
}

func Test_WithConfig(t *testing.T) {
	NewRuntime(WithConfig(&testConfig{})).Run(testApp{})
}

func Test_NoConfig(t *testing.T) {
	NewRuntime().Run(testApp{})
}
