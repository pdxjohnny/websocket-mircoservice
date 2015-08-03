package client

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Socket *websocket.Conn
  Recv func(message []byte)
}

func (ws *Conn) Connect(url_string string) (err error) {
	u, err := url.Parse(url_string)
	if err != nil {
		return err
	}

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}

	wsHeaders := http.Header{
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsSocket, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	if err != nil {
		fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
		return err
	}
	ws.Socket = wsSocket
	return nil
}

func (ws *Conn) Read() (err error) {
	for {
		_, message, err := ws.Socket.ReadMessage()
		if err != nil {
			return err
		}
    go ws.Recv(message)
	}
}

func (ws *Conn) Write(message []byte) (err error) {
  err = ws.Socket.WriteMessage(1, message)
  if err != nil {
    return err
  }
  return nil
}