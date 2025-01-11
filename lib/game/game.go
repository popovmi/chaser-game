package game

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

//go:generate msgp

type Game struct {
	Players       map[string]*Player `msg:"players"`
	PortalNetwork *PortalNetwork     `msg:"portalNetwork"`
	Bricks        []*Brick           `msg:"bricks"`

	Counter      atomic.Uint64 `msg:"-"`
	Mu           sync.Mutex    `msg:"-"`
	PreviousTick int64         `msg:"-"`
}

func NewGame() *Game {
	brickW, brickH := 200.0, 40.0

	pl1 := newPortalLink()
	pl2 := newPortalLink()
	p1 := newPortal(pl1.ID, 500, 500)
	p2 := newPortal(pl1.ID, FieldWidth-500, FieldHeight-500)
	p3 := newPortal(pl2.ID, 500, FieldHeight-500)
	p4 := newPortal(pl2.ID, FieldWidth-500, 500)
	pl1.PortalIDs = []string{p1.ID, p2.ID}
	pl2.PortalIDs = []string{p3.ID, p4.ID}
	pn := newPortalNetwork(
		map[string]*Portal{p1.ID: p1, p2.ID: p2, p3.ID: p3, p4.ID: p4},
		map[string]*PortalLink{pl1.ID: pl1, pl2.ID: pl2},
	)

	g := &Game{
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
				if p.ID != np.ID && np.Touching(p) {
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
	now := time.Now()
	dt := now.UnixMilli() - g.PreviousTick

	wallHits := make(map[string]bool)
	touches := make(map[string]bool)

	g.Mu.Lock()
	defer g.Mu.Unlock()

	for k1, p1 := range g.Players {
		if p1.Status == PlayerStatusDead {
			if time.Since(p1.DeadAt).Seconds() >= RespawnTime {
				g.RespawnPlayer(p1)
			}
		}
		g.PortalNetwork.TeleportTick(p1)
		p1.Tick(float64(dt)/1000, g.Players)
		for _, brick := range g.Bricks {
			if !p1.IsHooked && p1.Hook == nil {
				brick.CollideAndBounce(p1)
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
	}
	g.PreviousTick += dt
	return wallHits, touches
}

func (g *Game) RespawnPlayer(p *Player) {
	g.findFreeSpot(p)
	p.RespawnedAt = time.Now()
	p.HP = MaxHP
	p.Status = PlayerStatusAlive
}
