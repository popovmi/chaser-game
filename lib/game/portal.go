package game

import (
	"time"

	"wars/lib/vector"
)

//go:generate msgp

type Portal struct {
	Pos vector.Vector2D `msg:"pos"`
}

func newPortal(x, y float64) *Portal {
	return &Portal{Pos: vector.Vector2D{X: x, Y: y}}
}

type PortalLink struct {
	P1       *Portal              `msg:"p1"`
	P2       *Portal              `msg:"p2"`
	LastUsed map[string]time.Time `msg:"lu"`
}

func NewPortalLink(x1, y1, x2, y2 float64) *PortalLink {
	return &PortalLink{newPortal(x1, y1), newPortal(x2, y2), make(map[string]time.Time)}
}

func (p *Portal) Touching(plr *Player) bool {
	return p.Pos.Distance(plr.Position) <= (PortalRadius - Radius)
}

func (p *PortalLink) Touching(plr *Player) bool {
	return p.P1.Touching(plr) || p.P2.Touching(plr)
}

func (p *PortalLink) GetPortalUsage(plr *Player) (bool, time.Time) {
	touching := p.Touching(plr)
	if !touching {
		return false, time.Time{}
	}

	_, usedTimeAgo := p.IsLinkOnCooldown(plr)
	return touching, usedTimeAgo
}

func (p *PortalLink) IsLinkOnCooldown(plr *Player) (bool, time.Time) {
	if lu, used := p.LastUsed[plr.ID]; used {
		if time.Since(lu).Seconds() < PortalCooldown {
			return true, lu
		}
	}
	return false, time.Time{}
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
		p.LastUsed[plr.ID] = time.Now()
	}
	return ported
}

func (p *Portal) CollideAndTeleport(plr *Player, dest *Portal) bool {
	if !p.Touching(plr) {
		return false
	}

	plr.Position.X = dest.Pos.X
	plr.Position.Y = dest.Pos.Y

	return true
}
