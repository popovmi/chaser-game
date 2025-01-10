package vector

import "math"

//go:generate msgp

type Vector2D struct {
	X float64 `msg:"x"`
	Y float64 `msg:"y"`
}

func NewVector2D(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}

func (v *Vector2D) Add(x, y float64) {
	v.X += x
	v.Y += y
}

func (v *Vector2D) AddV(other Vector2D) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector2D) Sub(x, y float64) {
	v.X -= x
	v.Y -= y
}

func (v *Vector2D) SubV(other Vector2D) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vector2D) Mul(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector2D) Normalize() {
	length := v.Length()
	if length == 0 {
		return
	}
	v.X /= length
	v.Y /= length
}

func (v *Vector2D) LimitLength(maxLength float64) {
	length := v.Length()
	if length > maxLength {
		v.Normalize()
		v.Mul(maxLength)
	}
}

func (v *Vector2D) Distance(v2 Vector2D) float64 {
	dx := v2.X - v.X
	dy := v2.Y - v.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (v *Vector2D) Project(dir Vector2D) Vector2D {
	dirMagSq := dir.X*dir.X + dir.Y*dir.Y
	if dirMagSq < 1e-6 {
		return NewVector2D(0, 0)
	}
	newVector := NewVector2D(dir.X, dir.Y)
	newVector.Mul(v.DotProduct(dir) / dirMagSq)
	return newVector
}

func (v *Vector2D) DotProduct(v2 Vector2D) float64 {
	return v.X*v2.X + v.Y*v2.Y
}
