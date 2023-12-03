package bots

import (
	"server/game/constants"
	"server/game/object"
	"server/game/utils"
)

type Velocity struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Bot struct {
	State    *object.ObjectState `json:"state"`
	Velocity *Velocity
	ID       int `json:"id"`
}

func NewBot(id int) *Bot {
	state := object.NewObjectState(15, 15, *utils.GetRandomCoordinate())
	return &Bot{
		State: state,
		Velocity: &Velocity{
			X: 1,
			Y: 1,
		},
		ID: id,
	}
}

func (b *Bot) MoveRandomly() {
	if b.State.TopLeft.Y+b.Velocity.Y > constants.Boundaries["maxY"] || b.State.BottomLeft.Y+b.Velocity.Y < constants.Boundaries["minY"] {
		b.Velocity.Y *= -1
	}
	if b.State.TopRight.X+b.Velocity.X > constants.Boundaries["maxX"] || b.State.TopLeft.X+b.Velocity.X < constants.Boundaries["minX"] {
		b.Velocity.X *= -1
	}
	b.State.UpdateState(b.Velocity.X, b.Velocity.Y)
}
