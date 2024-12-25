package vector

import "math"

//go:generate msgp

type Vect2D struct {
	X float64 `msg:"x"`
	Y float64 `msg:"y"`
}

func NewVect2D(x, y float64) *Vect2D {
	return &Vect2D{X: x, Y: y}
}

func (v *Vect2D) Distance(v2 *Vect2D) float64 {
	dx := v2.X - v.X
	dy := v2.Y - v.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (v *Vect2D) Add(v2 *Vect2D) *Vect2D {
	return NewVect2D(v.X+v2.X, v.Y+v2.Y)
}

func (v *Vect2D) Subtract(v2 *Vect2D) *Vect2D {
	return NewVect2D(v.X-v2.X, v.Y-v2.Y)
}

func (v *Vect2D) Multiply(scalar float64) *Vect2D {
	return NewVect2D(v.X*scalar, v.Y*scalar)
}

func (v *Vect2D) Divide(scalar float64) *Vect2D {
	if scalar == 0 {
		panic("деление на ноль")
	}
	return NewVect2D(v.X/scalar, v.Y/scalar)
}

func (v *Vect2D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vect2D) Normalize() *Vect2D {
	mag := v.Magnitude()
	if mag > 0 {
		return NewVect2D(v.X/mag, v.Y/mag)
	}
	return v
}

func (v *Vect2D) DotProduct(v2 *Vect2D) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v *Vect2D) CrossProduct(v2 *Vect2D) float64 {
	return v.X*v2.Y - v.Y*v2.X
}

func (v *Vect2D) Project(dir *Vect2D) *Vect2D {
	dirMagSq := dir.X*dir.X + dir.Y*dir.Y
	if dirMagSq < 1e-6 { // Check for zero-length vector
		return NewVect2D(0, 0)
	}
	scalar := v.DotProduct(dir) / dirMagSq
	return dir.Multiply(scalar)
}

func (v *Vect2D) LimitMagnitude(maxValue float64) *Vect2D {
	mag := v.Magnitude()
	if mag <= maxValue {
		return v
	}
	return v.Normalize().Multiply(maxValue)
}

func (v *Vect2D) AngleBetween(v2 *Vect2D) float64 {
	dotProd := v.DotProduct(v2)
	magV1 := v.Magnitude()
	magV2 := v2.Magnitude()
	cosTheta := dotProd / (magV1 * magV2)
	return math.Acos(cosTheta)
}
