package ws

import (
	"fmt"
)

type Message struct {
	Content  string `json:"content"`
	Type     int    `json:"type"` //1 for announcement, 2 for users' messages
	PlayerID uint64 `json:"playerid"`
}

type WorldState struct {
	State []PlayerState `json:"state"`
}
type GamesManager struct {
	Games  map[uint64]*Game
	LastID uint64
}

type Game struct {
	Name    string             `json:"name"`
	ID      uint64             `json:"id"`
	Players map[string]*Player `json:"players"`
}

type Lobby struct {
	Chat       chan *Message
	WorldState chan []PlayerState
	Connect    chan *Player
	Disconnect chan *Player
	Players    map[uint64]*Player
	LastID     uint64
}

type Hub struct {
	GamesManager
	Lobby
}

func NewHub() *Hub {
	return &Hub{
		GamesManager: GamesManager{
			Games: make(map[uint64]*Game),
		},
		Lobby: Lobby{
			Chat:       make(chan *Message, 5),
			Players:    make(map[uint64]*Player),
			Connect:    make(chan *Player),
			Disconnect: make(chan *Player),
			WorldState: make(chan []PlayerState),
		},
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub goroutine running...")

	for {
		select {
		case globalM := <-h.Lobby.Chat:
			for _, player := range h.Lobby.Players {
				player.message <- globalM
			}
		case worldState := <-h.WorldState:
			for _, player := range h.Lobby.Players {
				player.worldState <- &WorldState{State: worldState}
			}
		case player := <-h.Lobby.Connect:
			h.Lobby.Players[player.ID] = player
			player.message <- &Message{
				Content:  "Welcome",
				PlayerID: player.ID,
			}
			fmt.Println("Added new player to lobby", h.Lobby.Players)
		case player := <-h.Lobby.Disconnect:
			delete(h.Lobby.Players, player.ID)
			h.SendToAllPlayer(&Message{
				Content:  "Disconnect",
				PlayerID: player.ID,
			})
			fmt.Println("Deleted player from lobby", h.Lobby.Players)
		}
	}
}

func (h *Hub) SendToAllPlayer(m *Message) {
	for _, p := range h.Players {
		p.message <- m
	}
}
