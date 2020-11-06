package pkg

import (
	"context"

	"github.com/vseinstrumentiru/lego/v2/internal/deprecated"
	"github.com/vseinstrumentiru/lego/v2/server"
)

// Deprecated: use server.Run(app, config)
func Run(ctx context.Context, app deprecated.App, opts ...server.Option) {
	wrappedApp, cfg := deprecated.NewApp(app)
	opts = append(opts, server.ConfigOption(cfg))
	server.Run(wrappedApp, opts...)
}
