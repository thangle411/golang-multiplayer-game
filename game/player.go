package game

import (
	"encoding/json"
	"log"
	"server/game/constants"
	"server/game/object"
	"server/game/point"

	"github.com/gorilla/websocket"
)

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

	return &Player{id, make(chan *Message, 5), make(chan State), ws, "", object.NewObjectState(10, 10, point.Point{X: 0, Y: 0}), 0}
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
		if player.State.TopLeft.Y+5 > constants.Boundaries["maxY"] {
			return
		}
		player.State.UpdateState(0, 3)
	case "arrow-up":
		if player.State.BottomLeft.Y-5 < constants.Boundaries["minY"] {
			return
		}
		player.State.UpdateState(0, -3)
	case "arrow-left":
		if player.State.BottomLeft.X-5 < constants.Boundaries["minX"] {
			return
		}
		player.State.UpdateState(-3, 0)
	case "arrow-right":
		if player.State.BottomRight.X+5 > constants.Boundaries["maxX"] {
			return
		}
		player.State.UpdateState(3, 0)
	}
}
