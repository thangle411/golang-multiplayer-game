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

// -- Path: /ws/joinGame/:gameID?clientID&clientName=name
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
