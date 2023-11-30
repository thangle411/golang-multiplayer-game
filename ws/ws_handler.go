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

// --- /createRoom
func (h *Handler) CreateRoom(c *gin.Context) {
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

	fmt.Println("New game being created: ", game)
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
	h.hub.Lobby.Join <- client
	go client.writeMessageTo()
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
