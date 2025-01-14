package game

import (
	"time"

	"github.com/matoous/go-nanoid/v2"

	"wars/lib/vector"
)

//go:generate msgp

type Portal struct {
	ID         string          `msg:"id"`
	LinkID     string          `msg:"link_id,omitempty"`
	LastUsedAt time.Time       `msg:"last_used_at"`
	Pos        vector.Vector2D `msg:"pos,omitempty"`
}

func newPortal(linkID string, x, y float64) *Portal {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return &Portal{ID: id, LinkID: linkID, Pos: vector.Vector2D{X: x, Y: y}}
}

func (p *Portal) Touching(plr *Player) bool {
	return p.Pos.Distance(plr.Position) <= (PortalRadius - Radius)
}

type PortalLink struct {
	ID        string               `msg:"id"`
	PortalIDs []string             `msg:"portals,allownil"`
	LastUsed  map[string]time.Time `msg:"lu,omitempty"`
}

func newPortalLink() *PortalLink {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return &PortalLink{ID: id, LastUsed: make(map[string]time.Time)}
}

type PortalNetwork struct {
	Portals map[string]*Portal     `msg:"portals"`
	Links   map[string]*PortalLink `msg:"portal_links"`
}

func newPortalNetwork(portals map[string]*Portal, links map[string]*PortalLink) *PortalNetwork {
	return &PortalNetwork{portals, links}
}

func (pn *PortalNetwork) CanUsePortal(player *Player) (bool, *Portal, *time.Duration) {
	if player.IsHooked {
		return false, nil, nil
	}
	var portal *Portal
	for _, plr := range pn.Portals {
		if plr.Touching(player) {
			portal = plr
			break
		}
	}
	if portal == nil {
		return false, nil, nil
	}
	link := pn.Links[portal.LinkID]
	lastUsed, ok := link.LastUsed[player.ID]
	if !ok {
		return true, portal, nil
	}
	cooldown := time.Since(lastUsed)
	if cooldown.Seconds() < PortalCooldown {
		return false, portal, &cooldown
	}
	return true, portal, nil
}

func (pn *PortalNetwork) Teleport(player *Player) bool {
	can, departure, _ := pn.CanUsePortal(player)
	if !can {
		return false
	}
	link := pn.Links[departure.LinkID]
	for _, arrivalID := range link.PortalIDs {
		if arrivalID != departure.ID {
			player.Teleporting = true
			player.DepPortalID = departure.ID
			player.ArrPortalID = arrivalID
			now := time.Now()
			player.TeleportedAt = now
			departure.LastUsedAt = now
			pn.Portals[arrivalID].LastUsedAt = now
			link.LastUsed[player.ID] = now
			return true
		}
	}
	return false
}

func (pn *PortalNetwork) TeleportTick(player *Player) {
	if !player.Teleporting {
		return
	}
	depPort := pn.Portals[player.DepPortalID]
	arrPort := pn.Portals[player.ArrPortalID]
	link := pn.Links[depPort.LinkID]
	progress := time.Since(link.LastUsed[player.ID]).Seconds() / TeleportDuration
	if progress >= 0.5 && !player.Teleported {
		dx := arrPort.Pos.X - player.Position.X
		dy := arrPort.Pos.Y - player.Position.Y
		player.Position.Add(dx, dy)
		if player.Hook != nil {
			player.Hook.End.Add(dx, dy)
		}
		player.Teleported = true
	}
	if progress >= 1 {
		player.Teleporting = false
		player.Teleported = false
		player.DepPortalID = ""
		player.ArrPortalID = ""
	}
}

func (pn *PortalNetwork) Short() *PortalNetwork {
	short := &PortalNetwork{Portals: make(map[string]*Portal), Links: make(map[string]*PortalLink)}
	for _, port := range pn.Portals {
		short.Portals[port.ID] = &Portal{
			ID:         port.ID,
			LastUsedAt: port.LastUsedAt,
		}
	}
	for _, link := range pn.Links {
		short.Links[link.ID] = &PortalLink{
			ID:       link.ID,
			LastUsed: link.LastUsed,
		}
	}
	return short
}
