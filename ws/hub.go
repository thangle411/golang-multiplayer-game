package ws

import "fmt"

type Game struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Hub struct {
	Games map[string]*Game
}

func NewHub() *Hub {
	return &Hub{
		Games: make(map[string]*Game),
	}
}

func Init(number int) {
	for i := 0; i < number; i++ {
		fmt.Println("Creating games")
	}
}

func Run() {

}
