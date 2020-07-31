package lecontext

import (
	"context"
	"github.com/vseinstrumentiru/lego/tools/lecontext"
)

// deprecated: use github.com/vseinstrumentiru/lego/pkg/tools/lecontext
func Extractor(ctx context.Context) map[string]interface{} {
	return contexttools.Extractor(ctx)
}
