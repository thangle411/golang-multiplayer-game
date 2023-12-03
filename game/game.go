package game

import (
	"errors"
	"fmt"
	"math"
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
	Point point.Point `json:"point"`
	ID    int         `json:"id"`
}

type Round struct {
	Level   int      `json:"level"`
	Squares []Square `json:"squares"`
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
	square := Square{
		Point: GetRandomCoordinate(),
		ID:    1,
	}
	g.Round = Round{
		Level:   1,
		Squares: []Square{square},
	}
	g.IsGameInProgress = true
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
			if IsPlayerInside(s.Point, *player.State.Point) {
				reqMet = append(reqMet, true)
			}
		}
	}

	if len(reqMet) == len(g.Round.Squares) && len(reqMet) > 0 {
		g.Round.Level += 1
		tempS := []Square{}
		min := math.Min(float64(g.Round.Level-1), float64(len(data)))
		for i := 0; i < int(min); i++ {
			tempS = append(tempS, Square{
				Point: GetRandomCoordinate(),
				ID:    i + 1,
			})
		}
		g.Round.Squares = tempS
	}

	return data
}

func IsPlayerInside(squareP point.Point, playerP point.Point) bool {
	minX := squareP.X - 20
	maxX := squareP.X + 20
	minY := squareP.Y - 20
	maxY := squareP.Y + 20
	if playerP.X >= minX && playerP.X <= maxX && playerP.Y >= minY && playerP.Y <= maxY {
		return true
	}
	return false
}

func GetRandomCoordinate() point.Point {
	return point.Point{
		X: utils.RandomCoordinate(Boundaries["minX"], Boundaries["maxX"]),
		Y: utils.RandomCoordinate(Boundaries["minY"], Boundaries["maxY"]),
	}
}
