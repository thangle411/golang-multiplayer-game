package router

import (
	"net/http"
	"server/ws"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Init(wsHandler *ws.Handler) {
	r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	r.GET("/ws/joinGame/:gameID", wsHandler.JoinGame)
}

func Start(addr string) error {
	return r.Run(addr)
}
