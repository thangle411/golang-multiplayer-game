package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// --- /createGame
func (h *Handler) CreateGame(c *gin.Context) {

	if len(h.hub.GamesManager.Games) >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Maximum amount of games reached",
		})
		return
	}

	type tmp struct {
		Name string `json:"name"`
	}
	var game tmp
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newid := h.hub.GamesManager.LastID + 1
	newRoom := Game{
		Name:    game.Name,
		ID:      newid,
		Clients: make(map[string]*Client),
	}
	fmt.Println("New game being created: ", newRoom)

	h.hub.GamesManager.LastID = newid
	h.hub.GamesManager.Games[newid] = &newRoom

	c.JSON(http.StatusOK, game)
}

// --- /getGames
func (h *Handler) GetGames(c *gin.Context) {
	type tmp struct {
		Name            string `json:"name"`
		ID              uint64 `json:"id"`
		NumberOfPlayers int    `json:"numberOfPlayers"`
	}
	games := make([]tmp, 0)

	for _, g := range h.hub.GamesManager.Games {
		games = append(games, tmp{
			Name:            g.Name,
			ID:              g.ID,
			NumberOfPlayers: len(g.Clients),
		})
	}

	c.JSON(http.StatusOK, games)
}

// --- /ws/joinLobby
func (h *Handler) JoinLobby(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newid := h.hub.Lobby.LastID + 1
	client := NewClient(newid, conn)
	fmt.Println("New client joins lobby", client)
	h.hub.Lobby.LastID = newid
	h.hub.Lobby.Chat <- &Message{
		Content: fmt.Sprintln("Random person ", newid, " joins the lobby"),
		Type:    1,
	}
	h.hub.Lobby.Connect <- client
	go client.writeMessageTo()
	client.readMessageFrom(h.hub)
}

// --- /ws/joinGame/:gameID?clientID&clientName=name
func (h *Handler) JoinGame(c *gin.Context) {
	_, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	gameID := c.Param("gameID")
	clientID := c.Query("clientID")
	clientName := c.Query("clientName")

	fmt.Println(gameID, clientID, clientName)
}
