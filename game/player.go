package game

import (
	"encoding/json"
	"log"
	"server/game/object"

	"github.com/gorilla/websocket"
)

var Boundaries = map[string]int{
	"minX": -400,
	"maxX": 400,
	"minY": -400,
	"maxY": 400,
}

type Input struct {
	Key string `json:"key"`
}

type SendData struct {
	WorldState []PlayerState `json:"worldState"`
	Gamestate  Round         `json:"gameState"`
}

type PlayerState struct {
	PlayerState object.ObjectState `json:"state"`
	ID          uint64             `json:"id"`
}

type Player struct {
	ID      uint64
	message chan *Message
	state   chan State
	wsConn  *websocket.Conn
	Name    string
	State   *object.ObjectState
	Room    int //0 meaning in lobby
}

func NewPlayer(id uint64, ws *websocket.Conn) *Player {
	if ws == nil {
		log.Fatal("ws cannot be nil")
	}

	return &Player{id, make(chan *Message, 5), make(chan State), ws, "", object.NewObjectState(), 0}
}

func (player *Player) ReadMessageFrom(hub *Hub) {
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
func (player *Player) WriteMessageTo() {
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

func (player *Player) WriteWorldStateTo() {
	defer func() {
		player.wsConn.Close()
	}()

	for {
		state, ok := <-player.state
		if !ok {
			return
		}
		player.wsConn.WriteJSON(state)
	}
}

func (player *Player) handleInput(input string) {
	switch input {
	case "arrow-down":
		if player.State.Point.Y+1 > Boundaries["maxY"] {
			return
		}
		player.State.Point.Y = player.State.Point.Y + 1
	case "arrow-up":
		if player.State.Point.Y-1 < Boundaries["minY"] {
			return
		}
		player.State.Point.Y = player.State.Point.Y - 1
	case "arrow-left":
		if player.State.Point.X-1 < Boundaries["minX"] {
			return
		}
		player.State.Point.X = player.State.Point.X - 1
	case "arrow-right":
		if player.State.Point.X+1 > Boundaries["maxX"] {
			return
		}
		player.State.Point.X = player.State.Point.X + 1
	}
}
