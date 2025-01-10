package game

import (
	"math"
	"sync"
	"time"

	"chaser/lib/colors"
	"chaser/lib/vector"
)

//go:generate msgp
type Direction byte

const (
	ZeroDir Direction = iota
	PositiveDir
	NegativeDir
)

type Player struct {
	ID    string      `msg:"id"`
	Name  string      `msg:"name"`
	Color colors.RGBA `msg:"color"`

	JoinedAt     time.Time `msg:"joined_at"`
	LastChasedAt time.Time `msg:"last_chased_at"`
	ChaseCount   int       `msg:"chase_count"`

	Position vector.Vector2D `msg:"position"`
	Velocity vector.Vector2D `msg:"velocity"`
	Angle    float64         `msg:"angle"`
	MoveDir  Direction       `msg:"move_dir"`
	TurnDir  Direction       `msg:"turn_dir"`
	Strafing bool            `msg:"strafing"`

	Hook       *Hook     `msg:"hook"`
	HookedAt   time.Time `msg:"hooked_at"`
	IsHooked   bool      `msg:"is_hooked"`
	CaughtByID string    `msg:"caught_by_id"`

	Blinking  bool      `msg:"blinking"`
	BlinkedAt time.Time `msg:"blinked_at"`
	Blinked   bool      `msg:"blinked"`

	mu sync.Mutex
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:       id,
		Position: vector.NewVector2D(0, 0),
		Velocity: vector.NewVector2D(0, 0),
		Angle:    0,
	}
}

func (p *Player) Tick(dt float64, players map[string]*Player) {
	if !p.IsHooked {
		p.Friction(dt)
		p.HookTick(dt, players)
		p.BlinkTick()
		p.Turn(dt)

		if p.Hook == nil || !p.Hook.Stuck {
			p.Accelerate()
			p.Step(dt)
		}
	}
}

func (p *Player) BlinkTick() {
	if p.Blinking {
		progress := time.Since(p.BlinkedAt).Seconds() / BlinkDuration
		if progress >= 0.5 && !p.Blinked {
			p.Position.Add(blinkDistance*math.Cos(p.Angle), blinkDistance*math.Sin(p.Angle))
			p.Blinked = true
		}
		if progress >= 1 {
			p.Blinking = false
			p.Blinked = false
		}
	}
}

func (p *Player) Turn(dt float64) {
	if !p.Strafing && p.TurnDir != ZeroDir {
		var angle float64
		if p.MoveDir == ZeroDir {
			angle = turnAngle * dt
		} else {
			angle = moveTurnAngle * dt
		}
		if p.TurnDir == NegativeDir {
			angle = -angle
		}
		p.Angle += angle
		p.Angle = math.Mod(p.Angle, 2*math.Pi)
	}
}

func (p *Player) Accelerate() {
	if p.MoveDir != ZeroDir {
		var ax, ay, maxV float64
		if p.Strafing {
			ax, ay = -strafeAcceleration*math.Sin(p.Angle), strafeAcceleration*math.Cos(p.Angle)
			maxV = maxStrafeVelocity
		} else {
			ax, ay = acceleration*math.Cos(p.Angle), acceleration*math.Sin(p.Angle)
			maxV = maxVelocity
		}
		if p.MoveDir == NegativeDir {
			ax, ay = -ax, -ay
		}
		previousVelocity := p.Velocity.Length()
		if previousVelocity > maxVelocity {
			maxV = maxCollideVelocity
		}
		p.Velocity.Add(ax, ay)
		p.Velocity.LimitLength(maxV)
	}
}

func (p *Player) Clamp() bool {
	hit := false
	if p.Position.X < Radius {
		p.Position.X = Radius
		p.Velocity.X *= -wallElasticity
		hit = true
	}
	if p.Position.X > FieldWidth-Radius {
		p.Position.X = FieldWidth - Radius
		p.Velocity.X *= -wallElasticity
		hit = true
	}
	if p.Position.Y < Radius {
		p.Position.Y = Radius
		p.Velocity.Y *= -wallElasticity
		hit = true
	}
	if p.Position.Y > FieldHeight-Radius {
		p.Position.Y = FieldHeight - Radius
		p.Velocity.Y *= -wallElasticity
		hit = true
	}
	if hit {
		p.Velocity.LimitLength(maxCollideVelocity)
	}
	return hit
}

func (p *Player) Step(dt float64) {
	p.Position.Add(p.Velocity.X*dt, p.Velocity.Y*dt)
}

func (p *Player) Friction(dt float64) {
	fr := friction
	if p.Strafing {
		fr = strafeFriction
	}
	p.Velocity.Mul(math.Exp(-fr * dt))
}

func (p *Player) HandleBlink() {
	if !p.Blinking && time.Since(p.BlinkedAt).Seconds() >= BlinkCooldown {
		p.Blinking = true
		p.BlinkedAt = time.Now()
	}
}

func (p *Player) HandleStrafe(strafe bool) {
	p.Strafing = strafe
}

func (p *Player) HandleMove(dir Direction) {
	p.MoveDir = dir
}

func (p *Player) HandleTurn(dir Direction) {
	p.TurnDir = dir
}

func (p *Player) Touching(p2 *Player) bool {
	return p.Position.Distance(p2.Position) < 2*Radius
}

func (p *Player) Touchable() bool {
	return time.Since(p.LastChasedAt).Seconds() >= untouchableTime &&
		time.Since(p.JoinedAt).Seconds() >= untouchableTime
}

func (p *Player) CollidePlayer(p2 *Player) {
	direction := vector.NewVector2D(p2.Position.X, p2.Position.Y)
	direction.SubV(p.Position)
	distance := direction.Length()
	direction.Normalize()

	v1n := p.Velocity.Project(direction)
	v2n := p2.Velocity.Project(direction)

	v1t := vector.NewVector2D(p.Velocity.X, p.Velocity.Y)
	v1t.SubV(v1n)
	v2t := vector.NewVector2D(p2.Velocity.X, p2.Velocity.Y)
	v2t.SubV(v2n)

	v1n, v2n = v2n, v1n

	v1n.AddV(v1t)
	v1n.Mul(1)
	v1n.LimitLength(maxCollideVelocity)

	v2n.AddV(v2t)
	v2n.Mul(1)
	v2n.LimitLength(maxCollideVelocity)

	p.Velocity = v1n
	p2.Velocity = v2n

	displacement := 2*Radius - distance
	if displacement > 0 {
		direction.Mul(displacement / 2)

		p.Position.SubV(direction)
		p2.Position.AddV(direction)

		p.Clamp()
		p2.Clamp()
	}
}
