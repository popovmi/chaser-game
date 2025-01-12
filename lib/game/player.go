package game

import (
	"math"
	"sync"
	"time"

	"wars/lib/colors"
	"wars/lib/vector"
)

//go:generate msgp

//msgp:replace Direction with:int
type Direction int

const (
	DirectionNone Direction = iota
	DirectionPositive
	DirectionNegative
)

//msgp:replace Direction with:int
type PlayerStatus int

const (
	PlayerStatusAlive PlayerStatus = iota
	PlayerStatusDead
)

type Player struct {
	ID    string      `msg:"id"`
	Name  string      `msg:"name"`
	Color colors.RGBA `msg:"color"`

	JoinedAt time.Time `msg:"joined_at"`

	Status PlayerStatus `msg:"status"`
	HP     float64      `msg:"hp"`

	Kills  int `msg:"kills"`
	Deaths int `msg:"deaths"`

	DeathPos    vector.Vector2D `msg:"death_pos"`
	DeadAt      time.Time       `msg:"dead_at"`
	RespawnedAt time.Time       `msg:"respawned_at"`

	Position    vector.Vector2D `msg:"position"`
	Velocity    vector.Vector2D `msg:"velocity"`
	Angle       float64         `msg:"angle"`
	MoveDir     string          `msg:"move_dir"`
	RotationDir Direction       `msg:"turn_dir"`
	Boosting    bool            `msg:"is_boosting"`

	Hook       *Hook     `msg:"hook"`
	HookedAt   time.Time `msg:"hooked_at"`
	IsHooked   bool      `msg:"is_hooked"`
	CaughtByID string    `msg:"caught_by_id"`

	Blinking  bool      `msg:"blinking"`
	BlinkedAt time.Time `msg:"blinked_at"`
	Blinked   bool      `msg:"blinked"`

	Teleporting  bool      `msg:"teleporting"`
	DepPortalID  string    `msg:"dep_portal_id"`
	ArrPortalID  string    `msg:"arr_portal_id"`
	Teleported   bool      `msg:"teleported"`
	TeleportedAt time.Time `msg:"teleported_at"`

	mu sync.Mutex
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:       id,
		Position: vector.NewVector2D(0, 0),
		Velocity: vector.NewVector2D(0, 0),
		Angle:    0,
		HP:       MaxHP,
		Status:   PlayerStatusAlive,
	}
}

func (p *Player) Tick(dt float64, players map[string]*Player) {
	if p.Status == PlayerStatusDead {
		return
	}
	p.BlinkTick()
	if !p.IsHooked && !p.Teleporting {
		p.Rotate(dt)
		p.HookTick(dt, players)

		if p.Hook == nil || !p.Hook.Stuck {
			p.Accelerate(dt)
			p.Step(dt)
		}
	}
}

func (p *Player) BlinkTick() {
	if p.Blinking {
		progress := time.Since(p.BlinkedAt).Seconds() / BlinkDuration
		if progress >= 0.5 && !p.Blinked {
			dx, dy := blinkDistance*math.Cos(p.Angle), blinkDistance*math.Sin(p.Angle)
			p.Position.Add(dx, dy)
			if p.Hook != nil {
				p.Hook.End.Add(dx, dy)
			}
			p.Blinked = true
		}
		if progress >= 1 {
			p.Blinking = false
			p.Blinked = false
		}
	}
}

func (p *Player) Rotate(dt float64) {
	if p.RotationDir != DirectionNone {
		var angle float64
		if p.MoveDir == "" && !p.Boosting {
			angle = turnAngle * dt
		} else {
			angle = moveTurnAngle * dt
		}
		if p.RotationDir == DirectionNegative {
			angle = -angle
		}

		p.Angle += angle
		p.Angle = math.Mod(p.Angle, 2*math.Pi)
	}
}

func (p *Player) Accelerate(dt float64) {
	if p.MoveDir == "" && !p.Boosting {
		p.Velocity.LimitLength(maxCollideVelocity)
		return
	}
	var angle float64
	switch p.MoveDir {
	case "u":
		angle = -math.Pi / 2
	case "d":
		angle = math.Pi / 2
	case "l":
		angle = math.Pi
	case "r":
		angle = 0
	case "ul":
		angle = -3 * math.Pi / 4
	case "ur":
		angle = -math.Pi / 4
	case "dl":
		angle = 3 * math.Pi / 4
	case "dr":
		angle = math.Pi / 4
	}
	if p.MoveDir != "" {
		p.Velocity.Add(acceleration*math.Cos(angle)*dt, acceleration*math.Sin(angle)*dt)
	}
	if p.Boosting {
		p.Velocity.Add(boostAcceleration*math.Cos(p.Angle)*dt, boostAcceleration*math.Sin(p.Angle)*dt)
	}
	newSpeed := p.Velocity.Length()
	if newSpeed > maxCollideVelocity {
		p.Velocity.LimitLength(maxCollideVelocity)
	} else if newSpeed > maxBoostVelocity || p.Boosting {
		p.Velocity.LimitLength(maxBoostVelocity)
	} else if newSpeed > maxVelocity {
		p.Velocity.LimitLength(maxVelocity)
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
	return hit
}

func (p *Player) Step(dt float64) {
	p.Position.Add(p.Velocity.X*dt, p.Velocity.Y*dt)
}

func (p *Player) HandleBlink() {
	if !p.Blinking && time.Since(p.BlinkedAt).Seconds() >= BlinkCooldown {
		p.Blinking = true
		p.BlinkedAt = time.Now()
	}
}

func (p *Player) HandleMove(dir string) {
	p.MoveDir = dir
}

func (p *Player) HandleBoost(boosting bool) {
	if !p.IsHooked && (p.Hook == nil || !p.Hook.Stuck) {
		p.Boosting = boosting
	}
}

func (p *Player) HandleRotate(dir Direction) {
	p.RotationDir = dir
}

func (p *Player) Touching(p2 *Player) bool {
	return p.Position.Distance(p2.Position) < 2*Radius
}

func (p *Player) Touchable() bool {
	return p.Status != PlayerStatusDead &&
		time.Since(p.JoinedAt).Seconds() >= untouchableTime &&
		time.Since(p.RespawnedAt).Seconds() >= untouchableTime
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

func (p *Player) Brake() {
	p.MoveDir = ""
	p.Velocity.Mul(Braking)
}

func (p *Player) die() {
	p.Status = PlayerStatusDead
	p.DeadAt = time.Now()
	p.DeathPos = vector.NewVector2D(p.Position.X, p.Position.Y)
	p.Deaths += 1
	p.IsHooked = false
	p.CaughtByID = ""
	p.Velocity = vector.NewVector2D(0, 0)
	p.Blinking = false
	p.Blinked = false
	p.Angle = 0
	p.Hook = nil
}
