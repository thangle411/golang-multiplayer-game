package object

import "server/game/point"

type ObjectState struct {
	Point *point.Point
}

func NewObjectState() *ObjectState {
	return &ObjectState{
		Point: point.NewPoint(),
	}
}
