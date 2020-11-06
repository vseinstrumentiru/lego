package sentry

import (
	"emperror.dev/emperror"
	"emperror.dev/emperror/httperr"
	"emperror.dev/emperror/utils/keyvals"
	"github.com/getsentry/raven-go"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type handler struct {
	level  logur.Level
	client *raven.Client
	stop   bool
}

func Handler(endpoint string, level logur.Level, stop bool) *handler {
	client, err := raven.New(endpoint)

	emperror.Panic(err)

	return &handler{
		level:  level,
		client: client,
		stop:   stop,
	}
}

func (h *handler) LevelEnabled(level logur.Level) bool {
	return h.level >= level
}

func (h *handler) Handle(msg multilog.Entry) {
	var message string
	var extra raven.Extra
	var interfaces []raven.Interface

	if err, ok := msg.(error); ok {
		// Get HTTP request (if any)
		if req, ok := httperr.HTTPRequest(err); ok {
			interfaces = append(interfaces, raven.NewHttp(req))
		}
		message = err.Error()

		extra = keyvals.ToMap(emperror.Context(err))

		interfaces = append(
			interfaces,
			raven.NewException(
				err,
				raven.GetOrNewStacktrace(emperror.ExposeStackTrace(err), 1, 3, h.client.IncludePaths()),
			),
		)
	} else {
		message = msg.Message()
		extra = msg.Fields()
		interfaces = append(interfaces, &raven.Message{
			Message: message,
			Params:  msg.Details(),
		})
	}

	packet := raven.NewPacketWithExtra(
		message,
		extra,
		interfaces...,
	)

	h.client.Capture(packet, nil)
}

func (h *handler) StopPropagation() bool {
	return h.stop
}
