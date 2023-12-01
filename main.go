package main

import (
	"server/game"
	"server/router"
	"server/ws"
)

func main() {
	hub := ws.NewHub()
	game := game.NewGame(hub)
	wsHandler := ws.NewHandler(hub)
	go hub.Run()
	go game.Run()
	router.Init(wsHandler)
	router.Start("0.0.0.0:8080")
}
