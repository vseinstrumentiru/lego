package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	le "github.com/vseinstrumentiru/lego/container"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/internal/config/env"
	"github.com/vseinstrumentiru/lego/internal/container"
)

type testDefaults struct{}

func (testDefaults) SetDefaults(e config.Env) {}

type testValidates struct{}

func (testValidates) Validate() error { return nil }

type testDefaultsValidates struct{}

func (testDefaultsValidates) SetDefaults(e config.Env) {}
func (testDefaultsValidates) Validate() error          { return nil }

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

	res := parse(&test, "app")
	assert.NotEmpty(t, res.defaults)
	assert.NotNil(t, res.defaults["app.a1"])
	assert.NotNil(t, res.defaults["app.a3"])
	assert.NotNil(t, res.defaults["app.a5"])
	assert.NotNil(t, res.defaults["app.a5.b1"])
	assert.NotEmpty(t, res.validates)
	assert.NotNil(t, res.validates["app.a2"])
	assert.NotNil(t, res.validates["app.a3"])
	assert.NotNil(t, res.validates["app.a5"])
	assert.NotNil(t, res.validates["app.a5.b2"])
	assert.NotNil(t, res.validates["app.a5.b2.c2"])
}

type TestRootConfig struct {
	A string
	B int
	C TestSub1Config
	D TestSub2Config
	E string
	F *TestSub1Config
}

type TestSub1Config struct {
	A string
	B int
}

type TestSub2Config struct {
	A int
	B TestSubSubConfig
	C string
}

type TestSubSubConfig struct {
	A string
	B int
	C string
}

func preloadEnvs() {
	_ = os.Setenv("TEST_A", "test-a")
	_ = os.Setenv("TEST_B", "1")
	_ = os.Setenv("TEST_C_A", "test-c-a")
	_ = os.Setenv("TEST_D_A", "2")
	_ = os.Setenv("TEST_D_B_A", "test-d-b-a")
	_ = os.Setenv("TEST_D_B_B", "3")
	_ = os.Setenv("TEST_F_A", "test-f-a")
	_ = os.Setenv("TEST_F_B", "4")
}

func Test_Provide(t *testing.T) {
	preloadEnvs()

	var config TestRootConfig
	err := Configure(argsIn{
		Config:    &config,
		Env:       env.New("test"),
		Container: container.New(),
	})

	ass := assert.New(t)
	ass.Nil(err)
	ass.Equal("test-a", config.A)
	ass.Equal(1, config.B)
	ass.Equal("test-c-a", config.C.A)
	ass.Equal(2, config.D.A)
	ass.Equal("test-d-b-a", config.D.B.A)
	ass.Equal(3, config.D.B.B)
	ass.Equal("test-c-a", config.C.A)
	ass.Equal("test-f-a", config.F.A)
	ass.Equal(4, config.F.B)
}

func Test_ConfigInContainer(t *testing.T) {
	preloadEnvs()

	var config TestRootConfig
	container := container.New()
	_ = Configure(argsIn{
		Config:    config,
		Env:       env.New("test"),
		Container: container,
	})
	ass := assert.New(t)

	container.Execute(func(cfg *TestRootConfig) {
		ass.NotNil(cfg)
		ass.Equal("test-a", cfg.A)
		ass.Equal(1, cfg.B)
		ass.Equal("test-c-a", cfg.C.A)
		ass.Equal(2, cfg.D.A)
		ass.Equal("test-d-b-a", cfg.D.B.A)
		ass.Equal(3, cfg.D.B.B)
		ass.Equal("test-c-a", cfg.C.A)
		ass.Equal("test-f-a", cfg.F.A)
		ass.Equal(4, cfg.F.B)
	})

	container.Execute(func(in struct {
		le.In
		Cfg *TestSub1Config `name:"cfg.c"`
	}) {
		ass.NotNil(in.Cfg)
		ass.Equal("test-c-a", in.Cfg.A)
	})

	container.Execute(func(cfg *TestSub2Config) {
		ass.NotNil(cfg)
		ass.Equal(2, cfg.A)
		ass.Equal("test-d-b-a", cfg.B.A)
		ass.Equal(3, cfg.B.B)
	})

	container.Execute(func(cfg *TestSubSubConfig) {
		ass.NotNil(cfg)
		ass.Equal("test-d-b-a", cfg.A)
		ass.Equal(3, cfg.B)
	})

	container.Execute(func(in struct {
		le.In
		Cfg *TestSub1Config `name:"cfg.f"`
	}) {
		ass.NotNil(in.Cfg)
		ass.Equal("test-f-a", in.Cfg.A)
		ass.Equal(4, in.Cfg.B)
	})
}
