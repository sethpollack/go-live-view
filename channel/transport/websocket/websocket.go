package websocket

import (
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/websocket"
	"github.com/sethpollack/go-live-view/channel"
)

var _ channel.Transport = (*wsTransport)(nil)

type wsTransport struct {
	path string
}

func New(path string) *wsTransport {
	return &wsTransport{
		path: path,
	}
}

func (t *wsTransport) Path() string {
	return t.path
}

func (t *wsTransport) Serve(handle func(channel.Conn), w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		c, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		defer c.Close()

		handle(&wsConn{conn: c})
	}
}

var _ channel.Conn = (*wsConn)(nil)

type wsConn struct {
	conn net.Conn
}

func (t *wsConn) ReadMessage() ([]byte, error) {
	data, _, err := wsutil.ReadClientData(t.conn)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *wsConn) WriteMessage(data []byte) error {
	return wsutil.WriteServerMessage(t.conn, ws.OpText, data)
}
