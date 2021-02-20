package netx

import (
	"net"

	"emperror.dev/errors"
	"github.com/cloudflare/tableflip"
)

func Listen(network, addr string, upg *tableflip.Upgrader) (net.Listener, error) {
	if upg != nil {
		return upg.Listen(network, addr)
	}

	ln, err := net.Listen(network, addr)

	return ln, errors.Wrap(err, "can't create new listener")
}
