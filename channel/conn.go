package channel

type Conn interface {
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
}

type conn struct {
	c Conn
}

func newConnection(c Conn) *conn {
	return &conn{c: c}
}

func (t *conn) ReadMessage() (*Message, error) {
	data, err := t.c.ReadMessage()
	if err != nil {
		return nil, err
	}
	return decode(data)
}

func (t *conn) WriteMessage(m *Message) error {
	data, err := encode(m)
	if err != nil {
		return err
	}
	return t.c.WriteMessage(data)
}
