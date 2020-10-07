package lego

import (
	"context"

	"emperror.dev/errors"

	"github.com/vseinstrumentiru/lego/server"
)

// deprecated: use server.Run(app, config)
func Run(ctx context.Context, app interface{}, opts ...interface{}) {
	if len(opts) == 0 {
		panic(errors.New("third argument must be pointer to config struct"))
	}

	server.Run(app, opts[0])
}
