package game

import (
	"log/slog"
	"math"
	"time"
)

//go:generate msgp

type Brick struct {
	Position *Vector `msg:"pos,omitempty" json:"Position,omitempty"`
	Width    float64 `msg:"w,omitempty" json:"Width,omitempty"`
	Height   float64 `msg:"h,omitempty" json:"Height,omitempty"`
	Angle    float64 `msg:"a,omitempty" json:"Angle,omitempty"`
}

func NewBrick(x, y, w, h, a float64) *Brick {
	return &Brick{NewVector(x, y), w, h, a}
}

type State struct {
	Players       map[string]*Player `msg:"ps,omitempty" json:"Players,omitempty"`
	PortalNetwork *PortalNetwork     `msg:"pn,omitempty" json:"PortalNetwork,omitempty"`
	Bricks        []*Brick           `msg:"br,omitempty" json:"Bricks,omitempty"`
}

func NewState() *State {
	brickW, brickH := float64(BrickWidth), float64(BrickHeight)

	pl1 := newPortalLink("1")
	pl2 := newPortalLink("2")
	p1 := newPortal("1", pl1.ID, 500, 500)
	p2 := newPortal("2", pl1.ID, FieldWidth-500, FieldHeight-500)
	p3 := newPortal("3", pl2.ID, 500, FieldHeight-500)
	p4 := newPortal("", pl2.ID, FieldWidth-500, 500)
	pl1.PortalIDs = []string{p1.ID, p2.ID}
	pl2.PortalIDs = []string{p3.ID, p4.ID}
	pn := newPortalNetwork(
		map[string]*Portal{p1.ID: p1, p2.ID: p2, p3.ID: p3, p4.ID: p4},

		map[string]*PortalLink{pl1.ID: pl1, pl2.ID: pl2},
	)

	return &State{
		Players:       make(map[string]*Player),
		PortalNetwork: pn,
		Bricks: []*Brick{

			// mid hor
			NewBrick(FieldWidth/2-brickW-Radius, FieldHeight/2-40-Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2-brickW-Radius, FieldHeight/2+Radius, brickW, brickH, math.Pi),
			NewBrick(FieldWidth/2+Radius, FieldHeight/2-brickH-Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2+Radius, FieldHeight/2+Radius, brickW, brickH, math.Pi),

			// mid ver
			NewBrick(
				FieldWidth/2-brickW/2-brickH/2-Radius,
				FieldHeight/2-brickW/2-brickH/2-brickH-3*Radius,
				brickW, brickH, math.Pi/2,
			),
			NewBrick(
				FieldWidth/2-brickW/2-brickH/2-Radius,
				FieldHeight/2+brickW/2+brickH/2+3*Radius,
				brickW, brickH, math.Pi/2,
			),
			NewBrick(
				FieldWidth/2-brickW/2+brickH/2+Radius,
				FieldHeight/2-brickW/2-brickH/2-brickH-3*Radius,
				brickW, brickH, math.Pi/2,
			),
			NewBrick(
				FieldWidth/2-brickW/2+brickH/2+Radius,
				FieldHeight/2+brickW/2+brickH/2+3*Radius,
				brickW, brickH, math.Pi/2,
			),

			// top hor
			NewBrick(FieldWidth/2-2*brickW-Radius, FieldHeight/2-brickW-2*brickH-5*Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2-brickW-Radius, FieldHeight/2-brickW-2*brickH-5*Radius, brickW, brickH, math.Pi),
			NewBrick(FieldWidth/2+Radius, FieldHeight/2-brickW-2*brickH-5*Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2+brickW+Radius, FieldHeight/2-brickW-2*brickH-5*Radius, brickW, brickH, math.Pi),

			// bottom hor
			NewBrick(FieldWidth/2-2*brickW-Radius, FieldHeight/2+brickW+brickH+5*Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2-brickW-Radius, FieldHeight/2+brickW+brickH+5*Radius, brickW, brickH, math.Pi),
			NewBrick(FieldWidth/2+Radius, FieldHeight/2+brickW+brickH+5*Radius, brickW, brickH, 0),
			NewBrick(FieldWidth/2+brickW+Radius, FieldHeight/2+brickW+brickH+5*Radius, brickW, brickH, math.Pi),

			NewBrick(FieldWidth/2-brickW/2, brickW/2-brickH/2+2*Radius, brickW, brickH, math.Pi/2),
			NewBrick(FieldWidth/2-brickW/2, FieldHeight-brickW/2-brickH/2-2*Radius, brickW, brickH, math.Pi/2),
		},
	}
}

func (g *Game) update(dt float64) {
	slog.Debug("updating state")
	s := g.State
	for k1, p1 := range s.Players {
		if p1.Status == PlayerStatusPreparing {
			continue
		}
		if p1.Status == PlayerStatusDead {
			if time.Since(*p1.DeadAt).Seconds() >= RespawnTime {
				g.spawn(p1)
				continue
			}
		}
		p1.mu.Lock()
		s.PortalNetwork.tick(p1)
		p1.tick(dt, s.Players)
		for _, brick := range s.Bricks {
			if p1.HookedBy == "" && p1.Hook == nil {
				p1.collideBrick(brick)
			}
		}
		for k2, p2 := range s.Players {
			if k1 < k2 && p1.Touchable() && p2.Touchable() {
				if p1.touchingPlayer(p2) {
					p2.mu.Lock()
					p1.collidePlayer(p2)
					p2.mu.Unlock()
				}
			}
		}
		p1.collideWall()
		p1.mu.Unlock()
	}
	slog.Debug("state updated", "time", time.Since(g.LastTick).Milliseconds())
}

func (g *Game) spawn(player *Player) {
	t := time.Now()
	player.mu.Lock()
	player.Status = PlayerStatusAlive
	player.SpawnedAt = &t
	player.HP = MaxHP
	player.Position = NewVector(FieldWidth/2, FieldHeight/2)
	player.Velocity = NewVector(0, 0)
	player.mu.Unlock()
	g.events <- Event{EventActionSpawned, player}
}
