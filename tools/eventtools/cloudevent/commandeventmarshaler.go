package cloudevent

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"time"
)

type CommandEventMarshaller struct {
	NewUUID      func() string
	GenerateName func(v interface{}) string
}

func (m CommandEventMarshaller) Marshal(v interface{}) (*message.Message, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(
		m.newUUID(),
		b,
	)
	msg.Metadata.Set(MetaName, m.Name(v))

	if timed, ok := v.(EventWithTime); ok {
		msg.Metadata[MetaCreatedAt] = timed.GetCreatedAt().Format(time.RFC3339)
	} else {
		msg.Metadata[MetaCreatedAt] = time.Now().Format(time.RFC3339)
	}

	return msg, nil
}

func (m CommandEventMarshaller) newUUID() string {
	if m.NewUUID != nil {
		return m.NewUUID()
	}

	// default
	return watermill.NewUUID()
}

func (CommandEventMarshaller) Unmarshal(msg *message.Message, v interface{}) (err error) {
	err = json.Unmarshal(msg.Payload, v)

	var timeStr string
	{
		var ok bool
		if timeStr, ok = msg.Metadata[MetaCreatedAt]; !ok {
			return errors.New("message not implemented CloudEvent spec")
		}
	}

	if timed, ok := v.(EventWithTime); ok && err == nil {
		t, err := time.Parse(time.RFC3339, timeStr)

		if err != nil {
			t, err = time.Parse(time.RFC3339Nano, timeStr)

			if err != nil {
				return err
			}
		}

		timed.SetCreatedAt(t)
	}

	return err
}

func (m CommandEventMarshaller) Name(cmdOrEvent interface{}) string {
	if m.GenerateName != nil {
		return m.GenerateName(cmdOrEvent)
	}

	return cqrs.FullyQualifiedStructName(cmdOrEvent)
}

func (m CommandEventMarshaller) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get(MetaName)
}
