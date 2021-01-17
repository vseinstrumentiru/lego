package provide

import (
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
)

func Pipeline() []interface{} {
	return []interface{}{
		func() (*tableflip.Upgrader, error) {
			return tableflip.New(tableflip.Options{})
		},
		func() *run.Group {
			return new(run.Group)
		},
	}
}
