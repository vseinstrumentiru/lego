package uprader

import (
	"github.com/cloudflare/tableflip"
)

func Provide() (*tableflip.Upgrader, error) {
	upg, err := tableflip.New(tableflip.Options{})

	return upg, err
}
