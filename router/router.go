package router

import (
	"net/http"
	"server/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Init(wsHandler *ws.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})
	r.GET("/getGames", wsHandler.GetGames)
	r.POST("/createUser")
	r.POST("/createGame", wsHandler.CreateGame)

	r.GET("/ws/joinLobby", wsHandler.JoinLobby)
	r.GET("/ws/joinGame/:gameID", wsHandler.JoinGame)
}

func Start(addr string) error {
	return r.Run(addr)
}
