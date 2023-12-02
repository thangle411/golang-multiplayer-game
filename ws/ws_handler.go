package ws

import (
	"fmt"
	"net/http"
	"server/game"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *game.Hub
}

func NewHandler(hub *game.Hub) *Handler {
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

// --- /ws/joinLobby
func (h *Handler) JoinLobby(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	current := h.hub.Lobby.LastID
	player := game.NewPlayer(current, conn)
	fmt.Println("New player joins lobby", player)
	h.hub.Lobby.LastID = current + 1
	h.hub.Lobby.Chat <- &game.Message{
		Content: fmt.Sprintln("Random person ", current, " joins the lobby"),
		Type:    1,
	}
	h.hub.Lobby.Connect <- player
	go player.WriteMessageTo()
	go player.WriteWorldStateTo()
	player.ReadMessageFrom(h.hub)
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
