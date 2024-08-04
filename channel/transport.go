package channel

type Transport interface {
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
}

type transport struct {
	t Transport
}

func NewTransport(t Transport) *transport {
	return &transport{t: t}
}

func (t *transport) ReadMessage() (*Message, error) {
	data, err := t.t.ReadMessage()
	if err != nil {
		return nil, err
	}
	return decode(data)
}

func (t *transport) WriteMessage(m *Message) error {
	data, err := encode(m)
	if err != nil {
		return err
	}
	return t.t.WriteMessage(data)
}
