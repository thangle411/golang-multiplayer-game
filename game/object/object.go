package object

import "server/game/point"

type Rectangle struct {
	TopLeft     *point.Point `json:"topLeft"`
	TopRight    *point.Point `json:"topRight"`
	BottomLeft  *point.Point `json:"bottomLeft"`
	BottomRight *point.Point `json:"bottomRight"`
	Center      *point.Point `json:"center"`
	Width       int          `json:"width"`
	Height      int          `json:"height"`
}

type ObjectState struct {
	Rectangle
}

func NewObjectState(width int, height int, center point.Point) *ObjectState {
	return &ObjectState{
		Rectangle: Rectangle{
			Width:  width,
			Height: height,
			Center: &point.Point{
				X: center.X,
				Y: center.Y,
			},
			TopLeft: &point.Point{
				X: -width/2 + center.X,
				Y: -height/2 + center.Y,
			},
			TopRight: &point.Point{
				X: width/2 + center.X,
				Y: -height/2 + center.Y,
			},
			BottomRight: &point.Point{
				X: width/2 + center.X,
				Y: height/2 + center.Y,
			},
			BottomLeft: &point.Point{
				X: -width/2 + center.X,
				Y: height/2 + center.Y,
			},
		},
	}
}

func (o *ObjectState) UpdateState(x int, y int) {
	o.BottomLeft.X += x
	o.BottomLeft.Y += y

	o.BottomRight.X += x
	o.BottomRight.Y += y

	o.TopLeft.X += x
	o.TopLeft.Y += y

	o.TopRight.X += x
	o.TopRight.Y += y

	o.Center.X += x
	o.Center.Y += y
}

func (o *ObjectState) CollisionDetection(otherObj ObjectState) bool {
	horizontalCheck1 := otherObj.BottomLeft.X < o.BottomLeft.X+o.Width
	horizontalCheck2 := otherObj.BottomLeft.X+otherObj.Width > o.BottomLeft.X
	verticalCheck1 := otherObj.BottomLeft.Y < o.BottomLeft.Y+o.Height
	verticalCheck2 := otherObj.BottomLeft.Y+otherObj.Height > o.BottomLeft.Y
	if horizontalCheck1 && horizontalCheck2 && verticalCheck1 && verticalCheck2 {
		return true
	}
	return false
}
