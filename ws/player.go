package ws

import (
	"encoding/json"
	"log"
	"server/game/object"

	"github.com/gorilla/websocket"
)

type Input struct {
	Key string `json:"key"`
}

type PlayerState struct {
	PlayerState object.ObjectState `json:"state"`
	ID          uint64             `json:"id"`
}

type Player struct {
	ID         uint64
	message    chan *Message
	worldState chan *WorldState
	wsConn     *websocket.Conn
	Name       string
	State      *object.ObjectState
	Room       int32 //0 meaning in lobby
}

func NewPlayer(id uint64, ws *websocket.Conn) *Player {
	if ws == nil {
		log.Fatal("ws cannot be nil")
	}

	return &Player{id, make(chan *Message, 5), make(chan *WorldState), ws, "", object.NewObjectState(), 0}
}

func (player *Player) readMessageFrom(hub *Hub) {
	defer func() {
		hub.Lobby.Disconnect <- player
		player.wsConn.Close()
	}()

	for {
		_, m, err := player.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var input Input
		if err := json.Unmarshal(m, &input); err != nil {
			message := Message{
				Content:  string(m),
				Type:     2,
				PlayerID: player.ID,
			}
			hub.Lobby.Chat <- &message
		} else {
			player.handleInput(input.Key)
		}
	}
}

// This will be run with a goroutine whenever a player joins the lobby or a game
func (player *Player) writeMessageTo() {
	defer func() {
		player.wsConn.Close()
	}()

	for {
		message, ok := <-player.message
		if !ok {
			return
		}
		player.wsConn.WriteJSON(message)
	}
}

func (player *Player) writeWorldStateTo() {
	defer func() {
		player.wsConn.Close()
	}()

	for {
		state, ok := <-player.worldState
		if !ok {
			return
		}
		player.wsConn.WriteJSON(state)
	}
}

func (player *Player) handleInput(input string) {
	switch input {
	case "arrow-down":
		player.State.Point.Y = player.State.Point.Y + 1
	case "arrow-up":
		player.State.Point.Y = player.State.Point.Y - 1
	case "arrow-left":
		player.State.Point.X = player.State.Point.X - 1
	case "arrow-right":
		player.State.Point.X = player.State.Point.X + 1
	}
}
