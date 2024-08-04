package channel

import (
	"github.com/go-json-experiment/json"
)

type Message struct {
	JoinRef string `json:"join_ref"`
	Ref     string `json:"ref"`
	Topic   string `json:"topic"`
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}

func encode(m *Message) ([]byte, error) {
	return json.Marshal([]any{
		m.JoinRef,
		m.Ref,
		m.Topic,
		m.Event,
		m.Payload,
	}, json.DefaultOptionsV2())
}

func decode(payload []byte) (*Message, error) {
	arr := make([]any, 5)
	err := json.Unmarshal(payload, &arr)
	if err != nil {
		msg := decodeBinary(payload)
		return msg, nil
	}

	msg := &Message{}

	if v, ok := arr[0].(string); ok {
		msg.JoinRef = v
	}
	if v, ok := arr[1].(string); ok {
		msg.Ref = v
	}
	if v, ok := arr[2].(string); ok {
		msg.Topic = v
	}
	if v, ok := arr[3].(string); ok {
		msg.Event = v
	}

	msg.Payload = arr[4]

	return msg, nil
}

func decodeBinary(buffer []byte) *Message {
	joinRefSize := buffer[1]
	refSize := buffer[2]
	topicSize := buffer[3]
	eventSize := buffer[4]

	offset := 5

	joinRef := string(buffer[offset : offset+int(joinRefSize)])
	offset += int(joinRefSize)

	ref := string(buffer[offset : offset+int(refSize)])
	offset += int(refSize)

	topic := string(buffer[offset : offset+int(topicSize)])
	offset += int(topicSize)

	event := string(buffer[offset : offset+int(eventSize)])
	offset += int(eventSize)

	data := buffer[offset:]

	return &Message{
		JoinRef: joinRef,
		Ref:     ref,
		Topic:   topic,
		Event:   event,
		Payload: data,
	}
}
