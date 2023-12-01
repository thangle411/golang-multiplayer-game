package game

import (
	"fmt"
	"server/ws"
	"time"
)

type Game struct {
	Hub *ws.Hub
}

func NewGame(hub *ws.Hub) *Game {
	return &Game{
		hub,
	}
}

func (game *Game) Run() {
	fmt.Println("Running game...")
	for {
		data := []ws.PlayerState{}
		for _, player := range game.Hub.Players {
			data = append(data, ws.PlayerState{
				PlayerState: *player.State,
				ID:          player.ID,
			})
		}
		game.Hub.WorldState <- data
		time.Sleep(10 * time.Millisecond)
	}
}
