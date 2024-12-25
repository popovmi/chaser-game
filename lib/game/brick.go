package warsgame

import (
	"math"
	"wars/lib/vector"
)

//go:generate msgp

type Brick struct {
	Pos *vector.Vect2D `msg:"pos"`
	W   float64        `msg:"w"`
	H   float64        `msg:"h"`
	A   float64        `msg:"a"`
}

func NewBrick(x, y, w, h, a float64) *Brick {
	return &Brick{&vector.Vect2D{X: x, Y: y}, w, h, a}
}

func (b *Brick) CollideAndBounce(plr *Player) bool {
	cx := b.Pos.X + b.W/2
	cy := b.Pos.Y + b.H/2

	rotatedPlayerX := (plr.Pos.X-cx)*math.Cos(-b.A) - (plr.Pos.Y-cy)*math.Sin(-b.A) + cx
	rotatedPlayerY := (plr.Pos.X-cx)*math.Sin(-b.A) + (plr.Pos.Y-cy)*math.Cos(-b.A) + cy

	nextRotatedPlayerX := (plr.Pos.X+plr.Vel.X-cx)*math.Cos(-b.A) - (plr.Pos.Y+plr.Vel.Y-cy)*math.Sin(-b.A) + cx
	nextRotatedPlayerY := (plr.Pos.X+plr.Vel.X-cx)*math.Sin(-b.A) + (plr.Pos.Y+plr.Vel.Y-cy)*math.Cos(-b.A) + cy

	closestX := math.Max(b.Pos.X, math.Min(rotatedPlayerX, b.Pos.X+b.W))
	closestY := math.Max(b.Pos.Y, math.Min(rotatedPlayerY, b.Pos.Y+b.H))

	distance := math.Sqrt(math.Pow(nextRotatedPlayerX-closestX, 2) + math.Pow(nextRotatedPlayerY-closestY, 2))

	if distance >= Radius {
		return false
	}

	nx := nextRotatedPlayerX - closestX
	ny := nextRotatedPlayerY - closestY

	rotatedPlayerVx := plr.Vel.X*math.Cos(-b.A) - plr.Vel.Y*math.Sin(-b.A)
	rotatedPlayerVy := plr.Vel.X*math.Sin(-b.A) + plr.Vel.Y*math.Cos(-b.A)

	if math.Abs(nx) > math.Abs(ny) {
		rotatedPlayerVx *= -1 * BrickElasticity
		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx)
	} else if math.Abs(nx) < math.Abs(ny) {
		rotatedPlayerVy *= -1 * BrickElasticity
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	} else {
		rotatedPlayerVx *= -1 * BrickElasticity
		rotatedPlayerVy *= -1 * BrickElasticity

		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx)
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	}

	plr.Vel.X = rotatedPlayerVx*math.Cos(b.A) - rotatedPlayerVy*math.Sin(b.A)
	plr.Vel.Y = rotatedPlayerVx*math.Sin(b.A) + rotatedPlayerVy*math.Cos(b.A)
	plr.Pos.X = (rotatedPlayerX-cx)*math.Cos(b.A) - (rotatedPlayerY-cy)*math.Sin(b.A) + cx
	plr.Pos.Y = (rotatedPlayerX-cx)*math.Sin(b.A) + (rotatedPlayerY-cy)*math.Cos(b.A) + cy

	return true
}
