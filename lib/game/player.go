package warsgame

import (
	"log"
	"strings"
	"sync"
	"time"
	"wars/lib/vector"

	"github.com/matoous/go-nanoid/v2"

	"wars/lib/color"
)

//go:generate msgp

type Hook struct {
	End             *vector.Vect2D `msg:"end"`
	Vel             *vector.Vect2D `msg:"vel"`
	Distance        float64        `msg:"dist"`
	CurrentDistance float64        `msg:"cDist"`
	CaughtPlayerID  string         `msg:"cpId"`
	IsReturning     bool           `msg:"ir"`
	WallStucked     bool           `msg:"stuck"`
}

type Player struct {
	ID           string         `msg:"id"`
	Name         string         `msg:"name"`
	Color        color.RGBA     `msg:"clr"`
	JoinedAt     int64          `msg:"ja"`
	Pos          *vector.Vect2D `msg:"pos"`
	Vel          *vector.Vect2D `msg:"vel"`
	Direction    string         `msg:"dir"`
	ChaseCount   int            `msg:"cc"`
	LastChasedAt int64          `msg:"lca"`
	Blinking     bool           `msg:"bl"`
	BlinkedAt    int64          `msg:"bla"`
	Hook         *Hook          `msg:"hook"`
	HookedAt     int64          `msg:"ha"`
	CaughtByID   string         `msg:"cbId"`
	IsHooked     bool           `msg:"ish"`

	mu sync.Mutex
}

func NewPlayer() *Player {
	id, err := gonanoid.New()
	if err != nil {
		log.Fatal(err)
	}
	return &Player{ID: id, Pos: vector.NewVect2D(0, 0), Vel: vector.NewVect2D(0, 0)}
}

func (p *Player) Tick() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsHooked || (p.Hook != nil && p.Hook.WallStucked) {
		return
	}

	var dvx, dvy float64

	switch p.Direction {
	case "l":
		dvx -= Acceleration
	case "r":
		dvx += Acceleration
	case "u":
		dvy -= Acceleration
	case "d":
		dvy += Acceleration
	case "lu":
		dvx -= Acceleration
		dvy -= Acceleration
	case "ru":
		dvx += Acceleration
		dvy -= Acceleration
	case "ld":
		dvx -= Acceleration
		dvy += Acceleration
	case "rd":
		dvx += Acceleration
		dvy += Acceleration
	}

	p.applyFriction()

	acc := vector.NewVect2D(dvx, dvy)
	currentSpeed := p.Vel.Magnitude()

	if currentSpeed > MaxVelocity {
		p.Vel = p.Vel.Add(acc).LimitMagnitude(MaxCollideVelocity)
	} else {
		p.Vel = p.Vel.Add(acc).LimitMagnitude(MaxVelocity)
	}

	p.Pos = p.Pos.Add(p.Vel)

	if p.Blinking {
		p.blink()
	}

}

func (p *Player) applyFriction() {
	if p.Vel.Magnitude() <= 0 {
		return
	}

	if Friction >= p.Vel.Magnitude() {
		p.Vel = vector.NewVect2D(0, 0)
	} else {
		p.Vel = p.Vel.Subtract(p.Vel.Normalize().Multiply(Friction))
	}
}

func (p *Player) checkWallHit() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	hit := false

	if p.Pos.X < Left {
		p.Pos.X = Left
		p.Vel.X = -p.Vel.X * WallElasticity
		hit = true
	} else if p.Pos.X > Right {
		p.Pos.X = Right
		p.Vel.X = -p.Vel.X * WallElasticity
		hit = true
	}

	if p.Pos.Y < Top {
		p.Pos.Y = Top
		p.Vel.Y = -p.Vel.Y * WallElasticity
		hit = true
	} else if p.Pos.Y > Bottom {
		p.Pos.Y = Bottom
		p.Vel.Y = -p.Vel.Y * WallElasticity
		hit = true
	}

	p.Vel = p.Vel.LimitMagnitude(MaxCollideVelocity)

	return hit
}

func (p *Player) ChangeDirection(dir string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Direction = dir
}

func (p *Player) Touching(p2 *Player) bool {
	return p.Pos.Distance(p2.Pos) < 2*Radius
}

func (p *Player) Touchable() bool {
	now := time.Now().UnixMilli()
	return now-p.LastChasedAt > UntouchableTime && now-p.JoinedAt > UntouchableTime
}

func (p *Player) Brake() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Direction = ""
	p.Vel = p.Vel.Multiply(Braking)
}

func (p *Player) blink() {
	p.Blinking = false
	if p.Vel.X == 0 && p.Vel.Y == 0 {
		return
	}
	p.Pos = p.Pos.Add(p.Vel.Normalize().Multiply(BlinkDistance))
	p.BlinkedAt = time.Now().UnixMilli()
}

func (p *Player) ThrowHook() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Hook != nil || time.Now().UnixMilli()-p.HookedAt < HookCooldown || (p.Vel.X == 0 && p.Vel.Y == 0) {
		return
	}

	p.Hook = &Hook{
		End:      vector.NewVect2D(p.Pos.X, p.Pos.Y),
		Vel:      p.Vel.Normalize().Multiply(HookVelocity),
		Distance: HookDistance,
	}
	p.HookedAt = time.Now().UnixMilli()
}

func (p *Player) Compare(ap *Player) int {
	return strings.Compare(p.ID, ap.ID)
}
