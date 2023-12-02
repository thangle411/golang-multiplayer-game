package game

import (
	"fmt"
)

type Message struct {
	Content  string `json:"content"`
	Type     int    `json:"type"` //1 for announcement, 2 for users' messages
	PlayerID uint64 `json:"playerid"`
}

type State struct {
	WorldState []PlayerState `json:"worldState"`
	GameState  Round         `json:"gameState"`
}

type Lobby struct {
	State      chan State
	Chat       chan *Message
	Connect    chan *Player
	Disconnect chan *Player
	Players    map[uint64]*Player
	LastID     uint64
}

type Hub struct {
	Lobby
}

func NewHub() *Hub {
	return &Hub{
		Lobby: Lobby{
			Chat:       make(chan *Message, 5),
			Players:    make(map[uint64]*Player),
			Connect:    make(chan *Player),
			Disconnect: make(chan *Player),
			State:      make(chan State),
		},
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub goroutine running...")

	for {
		select {
		case globalM := <-h.Lobby.Chat:
			for _, p := range h.Lobby.Players {
				p.message <- globalM
			}
		case state := <-h.State:
			for _, p := range h.Lobby.Players {
				p.state <- state
			}
		case p := <-h.Lobby.Connect:
			h.Lobby.Players[p.ID] = p
			p.message <- &Message{
				Content:  "Welcome",
				PlayerID: p.ID,
			}
			fmt.Println("Added new player to lobby", h.Lobby.Players)
		case p := <-h.Lobby.Disconnect:
			delete(h.Lobby.Players, p.ID)
			h.SendToAllPlayer(&Message{
				Content:  "Disconnect",
				PlayerID: p.ID,
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
