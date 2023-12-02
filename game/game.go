package game

import (
	"fmt"
	"server/game/point"
	"server/game/utils"
	"time"
)

type Game struct {
	Hub              *Hub
	isGameInProgress bool
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
		isGameInProgress: false,
	}
}

func (g *Game) Run() {
	fmt.Println("Running game...")
	for {
		//start game req
		if g.IsStandingOnStart() {
			g.StartGame()
		}

		data := g.CheckRoundStatus()

		g.Hub.State <- State{
			WorldState: data,
			GameState:  g.Round,
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (g *Game) StartGame() {
	fmt.Println("Starting game...")
	square := Square{
		Point: GetRandomCoordinate(),
		ID:    1,
	}
	g.Round = Round{
		Level:   1,
		Squares: []Square{square},
	}
	g.isGameInProgress = true
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

	fmt.Println(reqMet)
	fmt.Println(g.Round.Squares)

	if len(reqMet) == len(g.Round.Squares) && len(reqMet) > 0 {
		g.Round.Level += 1
		var tempS = []Square{}
		for i := 0; i < g.Round.Level-1; i++ {
			tempS = append(tempS, Square{
				Point: GetRandomCoordinate(),
				ID:    i + 1,
			})
		}
		g.Round.Squares = tempS
	}

	return data
}

func (g *Game) IsStandingOnStart() bool {
	if g.isGameInProgress {
		return false
	}
	for _, p := range g.Hub.Players {
		x := p.State.Point.X
		y := p.State.Point.Y
		if (x > -20 && x < 20) && (y > Boundaries["minY"] && y < Boundaries["minY"]+30) {
			return true
		}
	}
	return false
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
