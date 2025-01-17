package game

import "math"

//go:generate msgp

type Vector struct {
	X float64 `msg:"X"`
	Y float64 `msg:"Y"`
}

func NewVector(x, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

func (v *Vector) Clone() *Vector {
	return NewVector(v.X, v.Y)
}

func (v *Vector) Translate(dx, dy float64) {
	v.X += dx
	v.Y += dy
}

func (v *Vector) Add(other *Vector) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector) Subtract(other *Vector) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vector) Scale(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vector) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared())
}

func (v *Vector) Normalize() {
	magnitude := v.Magnitude()
	if magnitude != 0 {
		v.Scale(1 / magnitude)
	}
}

func (v *Vector) LimitMagnitude(maxMagnitude float64) {
	magnitudeSq := v.Magnitude()
	if magnitudeSq > maxMagnitude {
		v.Normalize()
		v.Scale(maxMagnitude)
	}
}

func (v *Vector) DistanceTo(other *Vector) float64 {
	return math.Sqrt(v.DistanceSquaredTo(other))
}

func (v *Vector) DistanceSquaredTo(other *Vector) float64 {
	dx := other.X - v.X
	dy := other.Y - v.Y
	return dx*dx + dy*dy
}

func (v *Vector) ProjectOnto(other *Vector) *Vector {
	dirMagSq := other.MagnitudeSquared()
	if dirMagSq < 1e-6 {
		return NewVector(0, 0)
	}
	scalarProjection := v.DotProduct(other) / dirMagSq
	return NewVector(other.X*scalarProjection, other.Y*scalarProjection)
}

func (v *Vector) DotProduct(other *Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}
