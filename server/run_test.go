package server

import "testing"

type testApp struct {
}

type testConfig struct {
}

func Test_Start(t *testing.T) {
	Run(testApp{}, &testConfig{}, NoWait)
}
