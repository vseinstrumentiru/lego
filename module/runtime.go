package module

import "github.com/vseinstrumentiru/lego/v2/module/pipeline"

func Pipeline() (interface{}, []interface{}) {
	return pipeline.Provide, []interface{}{
		pipeline.Configure,
	}
}
