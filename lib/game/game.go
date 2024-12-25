package warsgame

import (
	"math"
	"reflect"
	"sync"
	"time"
)

//go:generate msgp

type Game struct {
	Players     map[string]*Player `msg:"players"`
	CId         string             `msg:"cId"`
	PortalLinks []*PortalLink      `msg:"portalLinks"`
	Bricks      []*Brick           `msg:"bricks"`

	Mu sync.Mutex `msg:"-"`
}

func NewGame() *Game {
	brickW, brickH := 200.0, 40.0
	return &Game{
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
}

func (g *Game) SetPlayers(players map[string]*Player) {
	g.Players = players
}

func (g *Game) AddPlayer(p *Player) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.findFreeSpot(p)
	g.Players[p.ID] = p
	p.JoinedAt = time.Now().UnixMilli()

	if len(g.Players) == 1 {
		g.setChaser(p)
	}
}

func (g *Game) findFreeSpot(np *Player) {
	if len(g.Players) == 0 {
		np.Pos.X = FieldWidth / 2
		np.Pos.Y = FieldHeight / 2
		return
	}

	for y := FieldHeight / 2.0; y <= Bottom; y += 1 {
		for x := FieldWidth / 2.0; x <= Right; x += 1 {
			np.Pos.X = x
			np.Pos.Y = y
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
	g.Mu.Lock()
	defer g.Mu.Unlock()

	p := g.Players[id]
	delete(g.Players, id)
	if pV := reflect.ValueOf(p); !pV.IsNil() && p.ID == g.CId {
		if len(g.Players) == 0 {
			g.CId = ""
		} else {
			for _, p := range g.Players {
				g.setChaser(p)
				break
			}
		}
	}
}

func (g *Game) Tick() (map[string]bool, map[string]bool) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	l := len(g.Players)
	if l > 0 {
		for _, p := range g.Players {
			p.Tick()
			if p.Hook != nil {
				g.hookTick(p)
			}
		}
		return g.detectCollisions()
	}
	return make(map[string]bool), make(map[string]bool)
}

func (g *Game) Teleport(id string) bool {
	g.Mu.Lock()
	defer g.Mu.Unlock()

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

func (g *Game) CanUsePortal(id string) (bool, int64) {
	if p, ok := g.Players[id]; ok {
		if p.IsHooked {
			return false, 0
		}
		for _, link := range g.PortalLinks {
			touching, usedAt := link.GetPortalUsage(p)
			if touching {
				return touching, usedAt
			}
		}
	}
	return false, 0
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
		if p1.checkWallHit() {
			wallHits[p1.ID] = true
		}
		for k2, p2 := range g.Players {
			if k1 < k2 && p1.Touchable() && p2.Touchable() {
				if p1.Touching(p2) {
					touches[p1.ID] = true
					touches[p2.ID] = true

					direction := p2.Pos.Subtract(p1.Pos).Normalize()
					distance := direction.Magnitude()

					v1n := p1.Vel.Project(direction)
					v2n := p2.Vel.Project(direction)
					v1t := p1.Vel.Subtract(v1n)
					v2t := p2.Vel.Subtract(v2n)
					v1n, v2n = v2n, v1n

					p1.Vel = v1n.Add(v1t).Multiply(PlayerElasticity)
					p2.Vel = v2n.Add(v2t).Multiply(PlayerElasticity)
					p1.Vel.LimitMagnitude(MaxCollideVelocity)
					p2.Vel.LimitMagnitude(MaxCollideVelocity)

					displacement := Radius - distance
					if displacement > 0 {
						p1.Pos = p1.Pos.Subtract(direction.Multiply(displacement / 2))
						p2.Pos = p2.Pos.Add(direction.Multiply(displacement / 2))
					}

					if p1.ID == g.CId {
						g.setChaser(p2)
					} else if p2.ID == g.CId {
						g.setChaser(p1)
					}
				}
			}
		}
	}
	return wallHits, touches
}

func (g *Game) Blink(pId string) {
	if p, ok := g.Players[pId]; ok {
		if p.IsHooked || time.Now().UnixMilli()-p.BlinkedAt < BlinkCooldown {
			return
		}

		p.Blinking = true
	}

}

func (g *Game) hookTick(owner *Player) {
	h := owner.Hook

	if h.WallStucked {
		owner.Vel = h.End.Subtract(owner.Pos).Normalize().Multiply(HookBackwardVelocity)
		owner.Pos = owner.Pos.Add(owner.Vel)
		d := h.End.Distance(owner.Pos)
		if d <= Radius {
			owner.Hook = nil
		}
		return
	}

	if !h.IsReturning && h.CurrentDistance < h.Distance {
		h.End = h.End.Add(h.Vel)
		h.CurrentDistance += HookVelocity

		if h.End.X < 0 {
			h.End.X = 0
			h.WallStucked = true
		} else if h.End.X > FieldWidth {
			h.End.X = FieldWidth
			h.WallStucked = true
		}

		if h.End.Y < 0 {
			h.End.Y = 0
			h.WallStucked = true
		} else if h.End.Y > FieldHeight {
			h.End.Y = FieldHeight
			h.WallStucked = true
		}

		if h.WallStucked {
			return
		}

		for _, target := range g.Players {
			if target.ID == owner.ID {
				continue
			}

			distance := target.Pos.Distance(h.End)

			if distance < Radius {
				h.CaughtPlayerID = target.ID
				h.IsReturning = true
				target.IsHooked = true
				target.CaughtByID = owner.ID
				target.Vel.X = 0
				target.Vel.Y = 0
				target.Direction = ""
				break
			}
		}
	} else {
		h.IsReturning = true

		h.End = h.End.Add(owner.Pos.Subtract(h.End).Normalize().Multiply(HookBackwardVelocity))

		target, targetExists := g.Players[h.CaughtPlayerID]
		if targetExists {
			target.Pos.X = h.End.X
			target.Pos.Y = h.End.Y
		} else {
			for _, p := range g.Players {
				if p.ID == owner.ID {
					continue
				}

				if p.Pos.Subtract(h.End).Magnitude() < Radius && p.Touchable() {
					h.CaughtPlayerID = p.ID
					h.IsReturning = true
					p.IsHooked = true
					p.CaughtByID = owner.ID
					p.Vel.X = 0
					p.Vel.Y = 0
					p.Direction = ""
					targetExists = true
					target = p
					break
				}
			}
		}

		d := h.End.Distance(owner.Pos)
		if d < Radius {
			owner.Hook = nil
			if targetExists {
				target.IsHooked = false
				target.CaughtByID = ""
			}
		}
	}
}

func (g *Game) setChaser(p *Player) {
	if previousChaser, exists := g.Players[g.CId]; exists {
		previousChaser.LastChasedAt = time.Now().UnixMilli()
	}
	g.CId = p.ID
	p.ChaseCount++
}
