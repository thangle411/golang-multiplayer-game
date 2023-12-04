package point

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPoint(x int, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

func ResetPoint(p *Point) {
	p.X = 0
	p.Y = 0
}
