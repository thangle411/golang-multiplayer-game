package utils

import "math/rand"

func RandomCoordinate(min, max int) int {
	return rand.Intn(max-min+1) + min
}
