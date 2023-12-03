package utils

import (
	"math/rand"
	"server/game/constants"
	"server/game/point"
)

func RandomNum(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func GetRandomCoordinate() *point.Point {
	return &point.Point{
		X: RandomNum(constants.Boundaries["minX"], constants.Boundaries["maxX"]),
		Y: RandomNum(constants.Boundaries["minY"], constants.Boundaries["maxY"]),
	}
}
