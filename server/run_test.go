package server

import (
	"testing"

	"github.com/vseinstrumentiru/lego/v2/app"
)

type testApp struct {
}

type testConfig struct {
}

func Test_WithConfig(t *testing.T) {
	Run(testApp{}, app.WithConfig(&testConfig{}), app.NoWait())
}

func Test_NoConfig(t *testing.T) {
	Run(testApp{}, app.NoWait())
}
