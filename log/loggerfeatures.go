package log

import (
	"context"

	"logur.dev/logur"
)

func (l *logger) WithLevel(level logur.Level) Logger {
	if l.LevelEnabled(level) {
		return &logger{
			level:     level,
			handler:   l.handler,
			fields:    l.fields,
			extractor: l.extractor,
		}
	}

	return l
}

func (l *logger) WithHandler(handler EntryHandler) Logger {
	l.handler = AppendHandler(l.handler, handler)

	return l
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	if len(l.fields) > 0 {
		_fields := make(map[string]interface{}, len(l.fields)+len(fields))

		for key, value := range l.fields {
			_fields[key] = value
		}

		for key, value := range fields {
			_fields[key] = value
		}

		fields = _fields
	}

	return &logger{
		level:     l.level,
		handler:   l.handler,
		fields:    fields,
		extractor: l.extractor,
	}
}

func (l *logger) WithFilter(matcher EntryMatcher) Logger {
	return &logger{
		level:     l.level,
		handler:   WithFilter(l.handler, matcher),
		extractor: l.extractor,
	}
}

func (l *logger) WithErrFilter(matcher EntryErrMatcher) Logger {
	return &logger{
		level:     l.level,
		handler:   WithErrFilter(l.handler, matcher),
		extractor: l.extractor,
	}
}

func (l *logger) WithContext(ctx context.Context) Logger {
	if l.extractor == nil {
		return l
	}

	return l.WithFields(l.extractor(ctx))
}
