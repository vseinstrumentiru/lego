package module

import (
	"runtime"

	"github.com/vseinstrumentiru/lego/v2/module/uprader"
)

func Upgrader() (interface{}, []interface{}) {
	if runtime.GOOS == "windows" {
		return nil, nil
	}

	return uprader.Provide, []interface{}{uprader.Configure}
}
