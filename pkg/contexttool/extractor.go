package contexttool

import (
	"context"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	"go.opencensus.io/trace"
)

func Extractor(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if correlationID, ok := correlation.FromContext(ctx); ok {
		fields["correlation_id"] = correlationID
	}

	if operationName, ok := kitxendpoint.OperationName(ctx); ok {
		fields["operation_name"] = operationName
	}

	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()

		fields["trace_id"] = spanCtx.TraceID.String()
		fields["span_id"] = spanCtx.SpanID.String()
	}

	return fields
}
