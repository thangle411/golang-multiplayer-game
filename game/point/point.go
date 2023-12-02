package point

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPoint() *Point {
	return &Point{
		X: 0,
		Y: 0,
	}
}

func ResetPoint(p *Point) {
	p.X = 0
	p.Y = 0
}
