package ws

import "fmt"

type Game struct {
	Name    string             `json:"name"`
	ID      uint64             `json:"id"`
	Players map[string]*Player `json:"players"`
}

type Message struct {
	Content string `json:"content"`
	Type    int    `json:"type"` //1 for announcement, 2 for users' messages
}

type Lobby struct {
	Chat    chan *Message
	Join    chan *Client
	Leave   chan *Client
	Clients map[uint64]*Client
	LastID  uint64
}

type Hub struct {
	Games map[uint64]*Game
	Lobby
}

func NewHub() *Hub {
	return &Hub{
		Games: make(map[uint64]*Game),
		Lobby: Lobby{
			Chat:    make(chan *Message, 5),
			Clients: make(map[uint64]*Client),
			Join:    make(chan *Client),
			Leave:   make(chan *Client),
		},
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub goroutine running...")
	for {
		select {
		case globalM := <-h.Lobby.Chat:
			for _, client := range h.Lobby.Clients {
				client.message <- globalM
			}
		case client := <-h.Lobby.Join:
			h.Lobby.Clients[client.id] = client
			fmt.Println("Added new client to lobby", h.Lobby.Clients)
		case client := <-h.Lobby.Leave:
			delete(h.Lobby.Clients, client.id)
			fmt.Println("Deleted client from lobby", h.Lobby.Clients)
		}
	}
}
