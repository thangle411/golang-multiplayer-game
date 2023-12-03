package game

import (
	"errors"
	"fmt"
	"server/game/bots"
	"server/game/object"
	"server/game/point"
	"server/game/utils"
	"time"
)

type Game struct {
	Hub              *Hub
	IsGameInProgress bool
	Round            Round
}

type Square struct {
	SquareState object.ObjectState `json:"state"`
	ID          int                `json:"id"`
}

type Round struct {
	Level   int         `json:"level"`
	Squares []Square    `json:"squares"`
	Bots    []*bots.Bot `json:"bots"`
}

func NewGame(hub *Hub) *Game {
	return &Game{
		Hub:              hub,
		IsGameInProgress: false,
	}
}

func (g *Game) Run() {
	fmt.Println("Running game...")
	for {
		data := g.CheckRoundStatus()

		g.Hub.State <- State{
			WorldState: data,
			GameState:  g.Round,
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (g *Game) StartGame() error {
	if g.IsGameInProgress {
		return errors.New("Game is running")
	}

	fmt.Println("Starting game...")
	g.IsGameInProgress = true
	g.Round.Level = 1
	g.createSquare(20)
	return nil
}

func (g *Game) EndGame() error {
	if !g.IsGameInProgress {
		return errors.New("Game is not running")
	}

	fmt.Println("Ending game...")
	g.Round = Round{
		Level:   0,
		Squares: []Square{},
		Bots:    []*bots.Bot{},
	}
	g.IsGameInProgress = false
	return nil
}

func (g *Game) CheckRoundStatus() []PlayerState {
	data := []PlayerState{}
	reqMet := make([]bool, 0)
	for _, player := range g.Hub.Players {
		data = append(data, PlayerState{
			PlayerState: *player.State,
			ID:          player.ID,
		})

		for _, s := range g.Round.Squares {
			if s.SquareState.CollisionDetection(*player.State) {
				reqMet = append(reqMet, true)
			}
		}
	}
	for _, b := range g.Round.Bots {
		b.MoveRandomly()
	}

	//reset squares count
	if len(g.Round.Squares) != len(data) {
		g.createSquare(20)
	}

	//next level
	if len(reqMet) == len(g.Round.Squares) && len(reqMet) > 0 {
		g.Round.Level += 1
		g.createSquare(20)
		g.AddBots(len(g.Hub.Players))
	}

	return data
}

func (g *Game) AddBots(num int) {
	fmt.Println(num)
	bArray := []*bots.Bot{}
	for i := 0; i < num; i++ {
		b := bots.NewBot(len(g.Round.Bots) + i)
		bArray = append(bArray, b)
	}
	g.Round.Bots = append(g.Round.Bots, bArray...)
}

func (g *Game) createSquare(size int) {
	if !g.IsGameInProgress {
		return
	}
	tempS := []Square{}
	for i := 0; i < len(g.Hub.Players); i++ {
		square := Square{
			SquareState: *object.NewObjectState(size, size, point.Point{X: 0, Y: 0}),
			ID:          i + 1,
		}
		rPoint := utils.GetRandomCoordinate()
		square.SquareState.UpdateState(rPoint.X, rPoint.Y)
		tempS = append(tempS, square)
	}
	g.Round.Squares = tempS
}
