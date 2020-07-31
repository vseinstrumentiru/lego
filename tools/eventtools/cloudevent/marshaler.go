package cloudevent

import (
	"bytes"
	"encoding/json"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
)

type cloudEventMsg struct {
	UUID            string          `json:"id"`
	Type            string          `json:"type"`
	Payload         json.RawMessage `json:"data"`
	Source          string          `json:"source"`
	SpecVersion     string          `json:"specversion"`
	DataContentType string          `json:"datacontenttype"`
	CreatedAt       string          `json:"time"`
}

const (
	MetaCreatedAt       = "created_at"
	MetaName            = "name"
	MetaSource          = "source"
	MetaDataContentType = "data_content_type"
)

type NamedEvent interface {
	EventName() string
}

func EventName(v interface{}) string {
	if e, ok := v.(NamedEvent); ok {
		return e.EventName()
	}

	return ""
}

type Marshaller struct {
	SpecVersion string
	Source      string
}

func (m Marshaller) Unmarshal(stanMsg *stan.Msg) (*message.Message, error) {
	buf := new(bytes.Buffer)

	_, err := buf.Write(stanMsg.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot write stan message data to buffer")
	}

	decoder := json.NewDecoder(buf)

	var decodedMsg cloudEventMsg
	if err := decoder.Decode(&decodedMsg); err != nil {
		return nil, errors.Wrap(err, "cannot decode message")
	}

	// creating clean message, to avoid invalid internal state with ack
	msg := message.NewMessage(decodedMsg.UUID, []byte(decodedMsg.Payload))
	msg.Metadata[MetaCreatedAt] = decodedMsg.CreatedAt
	msg.Metadata[MetaName] = decodedMsg.Type
	msg.Metadata[MetaDataContentType] = decodedMsg.DataContentType
	msg.Metadata[MetaSource] = decodedMsg.Source

	return msg, nil
}

func (m Marshaller) Marshal(topic string, msg *message.Message) ([]byte, error) {
	cloudMsg := cloudEventMsg{
		UUID:            msg.UUID,
		Type:            msg.Metadata[MetaName],
		Payload:         json.RawMessage(msg.Payload),
		Source:          msg.Metadata[MetaSource],
		SpecVersion:     m.SpecVersion,
		DataContentType: msg.Metadata[MetaDataContentType],
		CreatedAt:       msg.Metadata[MetaCreatedAt],
	}

	// todo - use pool
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	if err := encoder.Encode(cloudMsg); err != nil {
		return nil, errors.Wrap(err, "cannot encode message")
	}

	return buf.Bytes(), nil
}
