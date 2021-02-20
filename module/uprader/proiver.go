package uprader

import (
	"github.com/cloudflare/tableflip"
)

func Provide() (*tableflip.Upgrader, error) {
	println("++++++++++++++")
	upg, err := tableflip.New(tableflip.Options{})

	return upg, err
}
