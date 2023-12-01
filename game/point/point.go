package point

type Point struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
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
