package game

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

//go:generate msgp

type Game struct {
	Players     map[string]*Player `msg:"players"`
	PortalLinks []*PortalLink      `msg:"portalLinks"`
	Bricks      []*Brick           `msg:"bricks"`

	Counter      atomic.Uint64 `msg:"-"`
	Mu           sync.Mutex    `msg:"-"`
	PreviousTick int64         `msg:"-"`
}

func NewGame() *Game {
	brickW, brickH := 200.0, 40.0

	g := &Game{
		Players: make(map[string]*Player),
		PortalLinks: []*PortalLink{
			NewPortalLink(350, 350, FieldWidth-350, FieldHeight-350),
			NewPortalLink(350, FieldHeight-350, FieldWidth-350, 350),
		},
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

	return g
}

func (g *Game) AddPlayer(p *Player) {
	g.findFreeSpot(p)
	g.Players[p.ID] = p
	p.JoinedAt = time.Now()
}

func (g *Game) findFreeSpot(np *Player) {
	if len(g.Players) == 0 {
		np.Position.X = FieldWidth / 2
		np.Position.Y = FieldHeight / 2
		return
	}

	for y := FieldHeight / 2.0; y <= FieldHeight-Radius; y += 1 {
		for x := FieldWidth / 2.0; x <= FieldWidth-Radius; x += 1 {
			np.Position.X = x
			np.Position.Y = y
			intersects := false
			for _, p := range g.Players {
				if np.Touching(p) {
					intersects = true
					break
				}
			}
			if !intersects {
				return
			}
		}
	}
}

func (g *Game) RemovePlayer(id string) {
	delete(g.Players, id)
}

func (g *Game) Tick() (map[string]bool, map[string]bool) {
	now := time.Now().UnixMilli()
	dt := now - g.PreviousTick

	for _, p := range g.Players {
		p.Tick(float64(dt)/1000, g.Players)
	}

	wallHits, touches := g.detectCollisions()
	g.PreviousTick += dt
	return wallHits, touches
}

func (g *Game) detectCollisions() (map[string]bool, map[string]bool) {
	wallHits := make(map[string]bool)
	touches := make(map[string]bool)

	for k1, p1 := range g.Players {
		for _, brick := range g.Bricks {
			if !p1.IsHooked && p1.Hook == nil {
				brick.CollideAndBounce(p1)
			}
		}
		if p1.Clamp() {
			wallHits[p1.ID] = true
		}
		for k2, p2 := range g.Players {
			if k1 < k2 && p1.Touchable() && p2.Touchable() {
				if p1.Touching(p2) {
					touches[p1.ID] = true
					touches[p2.ID] = true

					p1.CollidePlayer(p2)
				}
			}
		}
	}
	return wallHits, touches
}

func (g *Game) Teleport(id string) bool {
	if p, ok := g.Players[id]; ok {
		if p.IsHooked {
			return false
		}

		for _, link := range g.PortalLinks {
			if link.CollideAndTeleport(p) {
				return true
			}
		}
	}
	return false
}

func (g *Game) CanUsePortal(id string) (bool, time.Time) {
	if p, ok := g.Players[id]; ok {
		if p.IsHooked {
			return false, time.Time{}
		}
		for _, link := range g.PortalLinks {
			touching, usedAt := link.GetPortalUsage(p)
			if touching {
				return touching, usedAt
			}
		}
	}
	return false, time.Time{}
}
