package marshallers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"

	"github.com/vseinstrumentiru/lego/events"
	"github.com/vseinstrumentiru/lego/events/eventcoders"
)

type StanEncoder struct {
	events.TransportEncoder
}

type StanDecoder struct {
	events.TransportDecoder
}

func (s StanEncoder) Marshal(topic string, msg *message.Message) ([]byte, error) {
	return s.Encode(topic, msg)
}

func (s StanDecoder) Unmarshal(msg *stan.Msg) (*message.Message, error) {
	headers := make(map[string]string)
	headers[eventcoders.SequenceIDHeader] = fmt.Sprint(msg.Sequence)
	hash := md5.Sum([]byte(msg.Subject + headers[eventcoders.SequenceIDHeader]))
	headers[eventcoders.UUIDHeader] = hex.EncodeToString(hash[:])
	headers[eventcoders.TopicHeader] = msg.Subject
	headers[eventcoders.RFC3339NanoHeader] = time.Unix(msg.Timestamp, 0).Format(time.RFC3339Nano)

	if msg.Redelivered {
		headers[eventcoders.RedeliveredHeader] = fmt.Sprint(msg.RedeliveryCount)
	}

	return s.Decode(msg.Data, headers)
}
