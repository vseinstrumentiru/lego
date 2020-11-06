package server

import "testing"

type testApp struct {
}

type testConfig struct {
}

func Test_WithConfig(t *testing.T) {
	Run(testApp{}, ConfigOption(&testConfig{}), NoWaitOption())
}

func Test_NoConfig(t *testing.T) {
	Run(testApp{}, NoWaitOption())
}
