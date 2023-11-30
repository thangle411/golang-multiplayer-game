package ws

import "fmt"

type Message struct {
	Content  string `json:"content"`
	Type     int    `json:"type"` //1 for announcement, 2 for users' messages
	ClientID uint64 `json:"clientid"`
}

type UseMessage struct {
	Content string `json:"content"`
}
type GamesManager struct {
	Games  map[uint64]*Game
	LastID uint64
}

type Game struct {
	Name    string             `json:"name"`
	ID      uint64             `json:"id"`
	Clients map[string]*Client `json:"clients"`
}

type Lobby struct {
	Chat       chan *Message
	Connect    chan *Client
	Disconnect chan *Client
	Clients    map[uint64]*Client
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
			Clients:    make(map[uint64]*Client),
			Connect:    make(chan *Client),
			Disconnect: make(chan *Client),
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
		case client := <-h.Lobby.Connect:
			h.Lobby.Clients[client.id] = client
			fmt.Println("Added new client to lobby", h.Lobby.Clients)
		case client := <-h.Lobby.Disconnect:
			delete(h.Lobby.Clients, client.id)
			fmt.Println("Deleted client from lobby", h.Lobby.Clients)
		}
	}
}
