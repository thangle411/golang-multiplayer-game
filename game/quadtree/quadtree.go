package quadtree

import (
	"server/game/constants"
	"server/game/object"
	"server/game/point"
)

type QuadTree struct {
	Capacity  int
	Aabb      *object.ObjectState
	Objects   []*object.ObjectState
	NorthWest *QuadTree
	NorthEast *QuadTree
	SouthWest *QuadTree
	SouthEast *QuadTree
}

func NewQuadTree(cap int) *QuadTree {
	return &QuadTree{
		Capacity: cap,
		Aabb:     object.NewObjectState(constants.Boundaries["maxX"]-constants.Boundaries["minX"], constants.Boundaries["maxY"]-constants.Boundaries["minY"], *point.NewPoint(0, 0)),
		Objects:  make([]*object.ObjectState, 0),
		// northWest: NewQuadTree(4),
		// northEast: NewQuadTree(4),
		// southWest: NewQuadTree(4),
		// southEast: NewQuadTree(4),
	}
}

func (q *QuadTree) Subdivide() {
	//	|-----|-----|
	//	|	NW	|  NE |
	//	|     |     |
	//	|-----|-----|
	//	|	SW	| SE  |
	//	|     |     |
	//	|-----|-----|
	halfH := q.Aabb.Height / 2
	halfW := q.Aabb.Width / 2
	q.NorthWest = NewQuadTree(q.Capacity)
	q.NorthEast = NewQuadTree(q.Capacity)
	q.SouthWest = NewQuadTree(q.Capacity)
	q.SouthEast = NewQuadTree(q.Capacity)
	q.NorthWest.Aabb = object.NewObjectState(halfW, halfH, *point.NewPoint(q.Aabb.Center.X-halfW/2, q.Aabb.Center.Y-halfH/2))
	q.NorthEast.Aabb = object.NewObjectState(halfW, halfH, *point.NewPoint(q.Aabb.Center.X+halfW/2, q.Aabb.Center.Y-halfH/2))
	q.SouthWest.Aabb = object.NewObjectState(halfW, halfH, *point.NewPoint(q.Aabb.Center.X-halfW/2, q.Aabb.Center.Y+halfH/2))
	q.SouthEast.Aabb = object.NewObjectState(halfW, halfH, *point.NewPoint(q.Aabb.Center.X+halfW/2, q.Aabb.Center.Y+halfH/2))
}

func (q *QuadTree) Insert(o *object.ObjectState) bool {
	if !q.Aabb.Contains(*o) {
		return false
	}

	if len(q.Objects) < q.Capacity && q.NorthWest == nil {
		q.Objects = append(q.Objects, o)
		return true
	}

	if q.NorthWest == nil {
		q.Subdivide()
	}

	if q.NorthWest.Insert(o) {
		return true
	}
	if q.NorthEast.Insert(o) {
		return true
	}
	if q.SouthWest.Insert(o) {
		return true
	}
	if q.SouthEast.Insert(o) {
		return true
	}

	return false
}

func (q *QuadTree) queryRange() {

}
