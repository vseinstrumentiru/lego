package inject

import (
	"go.uber.org/dig"
)

type Constructor interface{}
type Invocation interface{}
type Interface interface{}
type RegisterOption = dig.ProvideOption

type In = dig.In
type Out = dig.Out
