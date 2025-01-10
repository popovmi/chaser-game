package game

import (
	"math"

	"wars/lib/vector"
)

//go:generate msgp

type Brick struct {
	Pos vector.Vector2D `msg:"pos"`
	W   float64         `msg:"w"`
	H   float64         `msg:"h"`
	A   float64         `msg:"a"`
}

func NewBrick(x, y, w, h, a float64) *Brick {
	return &Brick{vector.Vector2D{X: x, Y: y}, w, h, a}
}

func (b *Brick) CollideAndBounce(plr *Player) bool {
	cx := b.Pos.X + b.W/2
	cy := b.Pos.Y + b.H/2

	rotatedPlayerX := (plr.Position.X-cx)*math.Cos(-b.A) - (plr.Position.Y-cy)*math.Sin(-b.A) + cx
	rotatedPlayerY := (plr.Position.X-cx)*math.Sin(-b.A) + (plr.Position.Y-cy)*math.Cos(-b.A) + cy

	closestX := math.Max(b.Pos.X, math.Min(rotatedPlayerX, b.Pos.X+b.W))
	closestY := math.Max(b.Pos.Y, math.Min(rotatedPlayerY, b.Pos.Y+b.H))

	distance := math.Sqrt(math.Pow(rotatedPlayerX-closestX, 2) + math.Pow(rotatedPlayerY-closestY, 2))

	if distance >= Radius {
		return false
	}

	nx := rotatedPlayerX - closestX
	ny := rotatedPlayerY - closestY

	rotatedPlayerVx := plr.Velocity.X*math.Cos(-b.A) - plr.Velocity.Y*math.Sin(-b.A)
	rotatedPlayerVy := plr.Velocity.X*math.Sin(-b.A) + plr.Velocity.Y*math.Cos(-b.A)

	if math.Abs(nx) > math.Abs(ny) {
		rotatedPlayerVx *= -BrickElasticity
		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx) // Исправлено
	} else if math.Abs(nx) < math.Abs(ny) {
		rotatedPlayerVy *= -BrickElasticity
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny) // Исправлено
	} else {
		rotatedPlayerVx *= -BrickElasticity
		rotatedPlayerVy *= -BrickElasticity
		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx) // Исправлено
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny) // Исправлено
	}

	plr.Velocity.X = rotatedPlayerVx*math.Cos(b.A) - rotatedPlayerVy*math.Sin(b.A)
	plr.Velocity.Y = rotatedPlayerVx*math.Sin(b.A) + rotatedPlayerVy*math.Cos(b.A)
	plr.Position.X = (rotatedPlayerX-cx)*math.Cos(b.A) - (rotatedPlayerY-cy)*math.Sin(b.A) + cx
	plr.Position.Y = (rotatedPlayerX-cx)*math.Sin(b.A) + (rotatedPlayerY-cy)*math.Cos(b.A) + cy

	return true
}
