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
		Players: make(map[string]*Player),
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
	games := make([]tmp, len(h.hub.GamesManager.Games))

	for _, g := range h.hub.GamesManager.Games {
		games = append(games, tmp{
			Name:            g.Name,
			ID:              g.ID,
			NumberOfPlayers: len(g.Players),
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

	current := h.hub.Lobby.LastID
	player := NewPlayer(current, conn)
	fmt.Println("New player joins lobby", player)
	h.hub.Lobby.LastID = current + 1
	h.hub.Lobby.Chat <- &Message{
		Content: fmt.Sprintln("Random person ", current, " joins the lobby"),
		Type:    1,
	}
	h.hub.Lobby.Connect <- player
	go player.writeMessageTo()
	go player.writeWorldStateTo()
	player.readMessageFrom(h.hub)
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
