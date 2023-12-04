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
			Center: point.NewPoint(
				center.X,
				center.Y,
			),
			TopLeft: point.NewPoint(
				-width/2+center.X,
				-height/2+center.Y,
			),
			TopRight: point.NewPoint(
				width/2+center.X,
				-height/2+center.Y,
			),
			BottomRight: point.NewPoint(
				width/2+center.X,
				height/2+center.Y,
			),
			BottomLeft: point.NewPoint(
				-width/2+center.X,
				height/2+center.Y,
			),
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

func (o *ObjectState) Contains(otherObj ObjectState) bool {
	//(centerX - w/2, centerY - h/2)																(centerX + w/2, centerY - h/2)
	// |--------------------------------------------------------------------|
	// |								 																										|
	// |												centerX, centerY 														|
	// |								 																										|
	// |--------------------------------------------------------------------|
	//(centerX - w/2, centerY + h/2)																(centerX + w/2, centerY + h/2)

	centerX := o.Center.X
	centerY := o.Center.Y
	w := o.Width
	h := o.Height

	objectX := otherObj.Center.X
	objectY := otherObj.Center.Y
	if objectX > centerX-w/2 && objectX < centerX+w/2 && objectY > centerY-h/2 && objectY < centerY+h/2 {
		return true
	}
	return false
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
