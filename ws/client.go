package ws

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id      uint64
	message chan *Message
	wsConn  *websocket.Conn
}

type Player struct {
	Client
	Name string
}

func NewClient(id uint64, ws *websocket.Conn) *Client {
	if ws == nil {
		log.Fatal("ws cannot be nil")
	}

	return &Client{id, make(chan *Message, 5), ws}
}

func (client *Client) readMessageFrom(hub *Hub) {
	defer func() {
		hub.Lobby.Leave <- client
		client.wsConn.Close()
	}()

	for {
		_, m, err := client.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println("message from client", m)
	}
}

// This will be run with a goroutine whenever a client joins the lobby or a game
func (client *Client) writeMessageTo() {
	defer func() {
		client.wsConn.Close()
	}()

	for {
		message, ok := <-client.message
		if !ok {
			return
		}

		client.wsConn.WriteJSON(message)
	}
}
