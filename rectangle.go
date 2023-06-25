package main

type Rectangle struct {
	A FloatPoint
	B FloatPoint
}

func (me Rectangle) Shrink(delta float64) Rectangle {
	return Rectangle{
		A: FloatPoint{
			X: me.A.X + delta,
			Y: me.A.Y + delta,
		},
		B: FloatPoint{
			X: me.B.X - delta,
			Y: me.B.Y - delta,
		},
	}
}

func (me Rectangle) GetWidth() float64 {
	return me.B.X - me.A.X
}

func (me Rectangle) GetHeight() float64 {
	return me.B.Y - me.A.Y
}

func (me Rectangle) CheckContainsPoint(a FloatPoint) bool {
	return me.A.X <= a.X && a.X <= me.B.X &&
		me.A.Y <= a.Y && a.Y <= me.B.Y
}

func (me Rectangle) CheckCollides(other Rectangle) bool {
	return me.CheckContainsPoint(other.A) || me.CheckContainsPoint(other.B) ||
		me.CheckContainsPoint(FloatPoint{X: other.A.X, Y: other.B.Y}) ||
		me.CheckContainsPoint(FloatPoint{X: other.B.X, Y: other.A.Y})
}
