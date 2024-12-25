package warsgame

import (
	"time"
	"wars/lib/vector"
)

//go:generate msgp

type Portal struct {
	Pos *vector.Vect2D `msg:"pos"`
}

func newPortal(x, y float64) *Portal {
	return &Portal{Pos: &vector.Vect2D{X: x, Y: y}}
}

type PortalLink struct {
	P1       *Portal          `msg:"p1"`
	P2       *Portal          `msg:"p2"`
	LastUsed map[string]int64 `msg:"LastUsed"`
}

func NewPortalLink(x1, y1, x2, y2 float64) *PortalLink {
	return &PortalLink{newPortal(x1, y1), newPortal(x2, y2), make(map[string]int64)}
}

func (p *Portal) Touching(plr *Player) bool {
	return p.Pos.Distance(plr.Pos) <= (PortalRadius - Radius)
}

func (p *PortalLink) Touching(plr *Player) bool {
	return p.P1.Touching(plr) || p.P2.Touching(plr)
}

func (p *PortalLink) GetPortalUsage(plr *Player) (bool, int64) {
	touching := p.Touching(plr)
	if !touching {
		return false, 0
	}

	_, usedTimeAgo := p.IsLinkOnCooldown(plr)
	return touching, usedTimeAgo
}

func (p *PortalLink) IsLinkOnCooldown(plr *Player) (bool, int64) {
	now := time.Now().UnixMilli()
	if lu, used := p.LastUsed[plr.ID]; used {
		if now-lu < PortalCooldown {
			return true, lu
		}
	}
	return false, 0
}

func (p *PortalLink) CollideAndTeleport(plr *Player) bool {
	isCd, _ := p.IsLinkOnCooldown(plr)
	if isCd {
		return false
	}
	ported := p.P1.CollideAndTeleport(plr, p.P2)
	if !ported {
		ported = p.P2.CollideAndTeleport(plr, p.P1)
	}
	if ported {
		p.LastUsed[plr.ID] = time.Now().UnixMilli()
	}
	return ported
}

func (p *Portal) CollideAndTeleport(plr *Player, dest *Portal) bool {
	if !p.Touching(plr) {
		return false
	}

	plr.Pos.X = dest.Pos.X
	plr.Pos.Y = dest.Pos.Y

	return true
}
