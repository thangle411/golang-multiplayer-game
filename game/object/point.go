package object

type Point struct {
	x int32
	y int32
}

func NewPoint() *Point {
	return &Point{
		x: 0,
		y: 0,
	}
}
