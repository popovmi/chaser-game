package game

import (
	"fmt"
	"math"
	"sync"
	"time"
)

//go:generate msgp

type Direction byte

const (
	DirectionNone Direction = iota
	DirectionPositive
	DirectionNegative
)

func (d *Direction) ExtensionType() int8 {
	return 100 // Выберите уникальный номер расширения
}

func (d *Direction) Len() int {
	return 1 // Размер Direction в байтах
}

func (d *Direction) MarshalBinaryTo(b []byte) error {
	b[0] = byte(*d)
	return nil
}

func (d *Direction) UnmarshalBinary(b []byte) error {
	*d = Direction(b[0])
	return nil
}

func (d Direction) String() string {
	switch d {
	case DirectionNone:
		return "None"
	case DirectionPositive:
		return "Positive"
	case DirectionNegative:
		return "Negative"
	default:
		return fmt.Sprintf("Unknown(%d)", byte(d))
	}
}

func (d Direction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *Direction) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1]

	switch str {
	case "None":
		*d = DirectionNone
	case "Positive":
		*d = DirectionPositive
	case "Negative":
		*d = DirectionNegative
	default:
		return fmt.Errorf("unknown direction: %s", str)

	}
	return nil
}

type PlayerStatus byte

const (
	PlayerStatusPreparing PlayerStatus = iota
	PlayerStatusAlive
	PlayerStatusDead
)

func (ps *PlayerStatus) ExtensionType() int8 {
	return 101 // Выберите уникальный номер расширения, отличный от Direction
}

func (ps *PlayerStatus) Len() int {
	return 1 // Размер PlayerStatus в байтах
}

func (ps *PlayerStatus) MarshalBinaryTo(b []byte) error {
	b[0] = byte(*ps)
	return nil
}

func (ps *PlayerStatus) UnmarshalBinary(b []byte) error {
	*ps = PlayerStatus(b[0])
	return nil
}

func (ps PlayerStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ps.String())), nil
}

func (ps *PlayerStatus) UnmarshalJSON(b []byte) error {
	str := string(b)
	// Убираем кавычки, если они есть
	str = str[1 : len(str)-1]
	switch str {
	case "Preparing":
		*ps = PlayerStatusPreparing
	case "Alive":
		*ps = PlayerStatusAlive
	case "Dead":
		*ps = PlayerStatusDead
	default:
		return fmt.Errorf("unknown player status: %s", str)
	}
	return nil
}

func (p PlayerStatus) String() string {
	switch p {
	case PlayerStatusPreparing:
		return "Preparing"
	case PlayerStatusAlive:
		return "Alive"
	case PlayerStatusDead:
		return "Dead"
	default:
		return fmt.Sprintf("Unknown(%d)", byte(p))
	}
}

type Player struct {
	ID                string       `msg:"id" json:"id"`
	Name              string       `msg:"n,omitempty" json:"Name,omitempty"`
	Color             *RGBA        `msg:"clr,omitempty" json:"Color,omitempty"`
	JoinedAt          *time.Time   `msg:"ja,omitempty" json:"JoinedAt,omitempty"`
	Status            PlayerStatus `msg:"s,omitempty" json:"Status,omitempty"`
	HP                float64      `msg:"hp,omitempty" json:"HP,omitempty"`
	Position          *Vector      `msg:"pos,omitempty" json:"Position,omitempty"`
	Velocity          *Vector      `msg:"vel,omitempty" json:"Velocity,omitempty"`
	Angle             float64      `msg:"ang,omitempty" json:"Angle,omitempty"`
	MoveDirection     string       `msg:"md,omitempty" json:"MoveDirection,omitempty"`
	RotationDirection Direction    `msg:"rd,omitempty" json:"RotationDirection,omitempty"`
	Boosting          bool         `msg:"b,omitempty" json:"Boosting,omitempty"`
	Kills             int          `msg:"kls,omitempty" json:"Kills,omitempty"`
	Deaths            int          `msg:"dts,omitempty" json:"Deaths,omitempty"`
	DeathPosition     *Vector      `msg:"dpos,omitempty" json:"DeathPosition,omitempty"`
	DeadAt            *time.Time   `msg:"da,omitempty" json:"DeadAt,omitempty"`
	SpawnedAt         *time.Time   `msg:"sa,omitempty" json:"SpawnedAt,omitempty"`
	Hook              *Hook        `msg:"hk,omitempty,allownil" json:"Hook,omitempty,allownil"`
	UsedHookAt        *time.Time   `msg:"hka,omitempty" json:"UsedHookAt,omitempty"`
	HookedBy          string       `msg:"hb,omitempty" json:"HookedBy,omitempty"`
	Blinking          bool         `msg:"bl,omitempty" json:"Blinking,omitempty"`
	BlinkedAt         *time.Time   `msg:"bla,omitempty" json:"BlinkedAt,omitempty"`
	Blinked           bool         `msg:"bld,omitempty" json:"Blinked,omitempty"`
	Teleporting       bool         `msg:"tlp,omitempty" json:"Teleporting,omitempty"`
	FromPortalID      string       `msg:"depp,omitempty" json:"FromPortalID,omitempty"`
	ToPortalID        string       `msg:"arrp,omitempty" json:"ToPortalID,omitempty"`
	Teleported        bool         `msg:"tld,omitempty" json:"Teleported,omitempty"`
	TeleportedAt      *time.Time   `msg:"tla,omitempty" json:"TeleportedAt,omitempty"`

	mu sync.Mutex
}

func NewPlayer() *Player {
	return &Player{
		ID:     generateID(),
		Status: PlayerStatusPreparing,
	}
}

func (p *Player) tick(dt float64, players map[string]*Player) {
	if p.HookedBy == "" && !p.Teleporting {
		p.blinkTick()
		p.rotate(dt)
		p.hookTick(dt, players)

		if p.Hook == nil || !p.Hook.Stuck {
			p.accelerate(dt)
			p.step(dt)
		}
	}
}

func (p *Player) blinkTick() {
	if p.Blinking {
		progress := time.Since(*p.BlinkedAt).Seconds() / BlinkDuration
		if progress >= 0.5 && !p.Blinked {
			dx, dy := blinkDistance*math.Cos(p.Angle), blinkDistance*math.Sin(p.Angle)
			p.translatePosition(dx, dy)
			p.Blinked = true
		}
		if progress >= 1 {
			p.Blinking = false
			p.Blinked = false
		}
	}
}

func (p *Player) accelerate(dt float64) {
	if p.MoveDirection == "" && !p.Boosting {
		p.Velocity.LimitMagnitude(maxCollideVelocity)
		return
	}
	var angle float64
	switch p.MoveDirection {
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
	if p.MoveDirection != "" {
		p.Velocity.Translate(acceleration*math.Cos(angle)*dt, acceleration*math.Sin(angle)*dt)
	}
	if p.Boosting {
		p.Velocity.Translate(boostAcceleration*math.Cos(p.Angle)*dt, boostAcceleration*math.Sin(p.Angle)*dt)
	}
	newSpeed := p.Velocity.Magnitude()
	if newSpeed > maxCollideVelocity {
		p.Velocity.LimitMagnitude(maxCollideVelocity)
	} else if newSpeed > maxBoostVelocity || p.Boosting {
		p.Velocity.LimitMagnitude(maxBoostVelocity)
	} else if newSpeed > maxVelocity {
		p.Velocity.LimitMagnitude(maxVelocity)
	}
}

func (p *Player) rotate(dt float64) {
	if p.RotationDirection != DirectionNone {
		var angle float64
		if p.MoveDirection == "" && !p.Boosting {
			angle = turnAngle * dt
		} else {
			angle = moveTurnAngle * dt
		}
		if p.RotationDirection == DirectionNegative {
			angle = -angle
		}

		p.Angle += angle
		p.Angle = math.Mod(p.Angle, 2*math.Pi)
	}
}

func (p *Player) step(dt float64) {
	p.Position.Translate(p.Velocity.X*dt, p.Velocity.Y*dt)
}

func (p *Player) translatePosition(dx, dy float64) {
	p.Position.Translate(dx, dy)
	if p.Hook != nil {
		p.Hook.EndPosition.Translate(dx, dy)
	}
}

func (p *Player) die() {
	t := time.Now()
	p.Status = PlayerStatusDead
	p.DeadAt = &t
	p.DeathPosition = NewVector(p.Position.X, p.Position.Y)
	p.Deaths += 1
	p.HookedBy = ""
	p.Velocity = NewVector(0, 0)
	p.Blinking = false
	p.Blinked = false
	p.Angle = 0
	p.Hook = nil
}

func (p *Player) Touchable() bool {
	return p.Status == PlayerStatusAlive &&
		time.Since(*p.JoinedAt).Seconds() >= untouchableTime &&
		time.Since(*p.SpawnedAt).Seconds() >= untouchableTime
}

func (p *Player) touchingPlayer(other *Player) bool {
	return p.Position.DistanceTo(other.Position) < 2*Radius
}

func (p *Player) touchingPortal(portal *Portal) bool {
	return p.Position.DistanceTo(portal.Position) < (PortalRadius - Radius)
}

func (p *Player) collideWall() bool {
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

func (p *Player) collidePlayer(p2 *Player) {
	direction := NewVector(p2.Position.X, p2.Position.Y)
	direction.Subtract(p.Position)
	distance := direction.Magnitude()
	direction.Normalize()

	v1n := p.Velocity.ProjectOnto(direction)
	v2n := p2.Velocity.ProjectOnto(direction)

	v1t := NewVector(p.Velocity.X, p.Velocity.Y)
	v1t.Subtract(v1n)
	v2t := NewVector(p2.Velocity.X, p2.Velocity.Y)
	v2t.Subtract(v2n)

	v1n, v2n = v2n, v1n

	v1n.Add(v1t)
	v1n.LimitMagnitude(maxCollideVelocity)

	v2n.Add(v2t)
	v2n.LimitMagnitude(maxCollideVelocity)

	p.Velocity = v1n
	p2.Velocity = v2n

	displacement := 2*Radius - distance
	if displacement > 0 {
		direction.Scale(displacement / 2)

		p.Position.Subtract(direction)
		p2.Position.Add(direction)
	}
}

func (p *Player) collideBrick(b *Brick) bool {
	cx := b.Position.X + b.Width/2
	cy := b.Position.Y + b.Height/2

	rotatedPlayerX := (p.Position.X-cx)*math.Cos(-b.Angle) - (p.Position.Y-cy)*math.Sin(-b.Angle) + cx
	rotatedPlayerY := (p.Position.X-cx)*math.Sin(-b.Angle) + (p.Position.Y-cy)*math.Cos(-b.Angle) + cy

	closestX := math.Max(b.Position.X, math.Min(rotatedPlayerX, b.Position.X+b.Width))
	closestY := math.Max(b.Position.Y, math.Min(rotatedPlayerY, b.Position.Y+b.Height))

	distance := math.Sqrt(math.Pow(rotatedPlayerX-closestX, 2) + math.Pow(rotatedPlayerY-closestY, 2))

	if distance >= Radius {
		return false
	}

	nx := rotatedPlayerX - closestX
	ny := rotatedPlayerY - closestY

	rotatedPlayerVx := p.Velocity.X*math.Cos(-b.Angle) - p.Velocity.Y*math.Sin(-b.Angle)
	rotatedPlayerVy := p.Velocity.X*math.Sin(-b.Angle) + p.Velocity.Y*math.Cos(-b.Angle)

	if math.Abs(nx) > math.Abs(ny) {
		rotatedPlayerVx *= -BrickElasticity
		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx)
	} else if math.Abs(nx) < math.Abs(ny) {
		rotatedPlayerVy *= -BrickElasticity
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	} else {
		rotatedPlayerVx *= -BrickElasticity
		rotatedPlayerVy *= -BrickElasticity
		rotatedPlayerX = closestX + (Radius+0.0001)*math.Copysign(1, nx)
		rotatedPlayerY = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	}

	p.Velocity.X = rotatedPlayerVx*math.Cos(b.Angle) - rotatedPlayerVy*math.Sin(b.Angle)
	p.Velocity.Y = rotatedPlayerVx*math.Sin(b.Angle) + rotatedPlayerVy*math.Cos(b.Angle)
	p.Position.X = (rotatedPlayerX-cx)*math.Cos(b.Angle) - (rotatedPlayerY-cy)*math.Sin(b.Angle) + cx
	p.Position.Y = (rotatedPlayerX-cx)*math.Sin(b.Angle) + (rotatedPlayerY-cy)*math.Cos(b.Angle) + cy

	return true
}

func (p *Player) Set(other *Player) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.JoinedAt = other.JoinedAt
	p.Status = other.Status
	p.HP = other.HP
	p.Position = other.Position
	p.Velocity = other.Velocity
	p.Angle = other.Angle
	p.MoveDirection = other.MoveDirection
	p.RotationDirection = other.RotationDirection
	p.Boosting = other.Boosting
	p.Kills = other.Kills
	p.Deaths = other.Deaths
	p.DeathPosition = other.DeathPosition
	p.DeadAt = other.DeadAt
	p.SpawnedAt = other.SpawnedAt
	p.Hook = other.Hook
	p.UsedHookAt = other.UsedHookAt
	p.HookedBy = other.HookedBy
	p.Blinking = other.Blinking
	p.BlinkedAt = other.BlinkedAt
	p.Blinked = other.Blinked
	p.Teleporting = other.Teleporting
	p.FromPortalID = other.FromPortalID
	p.ToPortalID = other.ToPortalID
	p.Teleported = other.Teleported
	p.TeleportedAt = other.TeleportedAt
}
