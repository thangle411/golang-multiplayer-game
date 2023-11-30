package main

import (
	"server/router"
	"server/ws"
)

func main() {
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	router.Init(wsHandler)
	router.Start("0.0.0.0:8080")
}
