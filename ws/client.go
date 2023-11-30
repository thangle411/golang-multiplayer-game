package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     uint32
	wsConn *websocket.Conn
}

func NewClient(ws *websocket.Conn) *Client {
	if ws == nil {
		log.Fatal("ws cannot be nil")
	}

	return &Client{0, ws}
}
