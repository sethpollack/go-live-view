package ws

import (
	"net"
	"net/http"

	"github.com/sethpollack/go-live-view/channel"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/websocket"
)

var _ channel.Transport = (*wsTransport)(nil)

type wsHandler struct {
	onStatic func(http.ResponseWriter, *http.Request)
	onWS     func(channel.Transport)
}

func NewWSHandler(
	onStatic func(http.ResponseWriter, *http.Request),
	onWS func(channel.Transport)) *wsHandler {
	return &wsHandler{onStatic, onWS}
}

func (h *wsHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		h.onStatic(w, r)
		return
	}

	if websocket.IsWebSocketUpgrade(r) {
		c, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		defer c.Close()

		h.onWS(&wsTransport{conn: c})
	}

}

type wsTransport struct {
	conn net.Conn
}

func (t *wsTransport) ReadMessage() ([]byte, error) {
	data, _, err := wsutil.ReadClientData(t.conn)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *wsTransport) WriteMessage(data []byte) error {
	return wsutil.WriteServerMessage(t.conn, ws.OpText, data)
}
