package game

import (
	"time"

	"github.com/matoous/go-nanoid/v2"

	"wars/lib/vector"
)

//go:generate msgp

type Portal struct {
	ID     string          `msg:"id"`
	LinkID string          `msg:"link_id"`
	Pos    vector.Vector2D `msg:"pos"`
}

func newPortal(linkID string, x, y float64) *Portal {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return &Portal{id, linkID, vector.Vector2D{X: x, Y: y}}
}

func (p *Portal) Touching(plr *Player) bool {
	return p.Pos.Distance(plr.Position) <= (PortalRadius - Radius)
}

type PortalLink struct {
	ID        string               `msg:"id"`
	PortalIDs []string             `msg:"portals"`
	LastUsed  map[string]time.Time `msg:"lu"`
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
			link.LastUsed[player.ID] = time.Now()
			dx := pn.Portals[arrivalID].Pos.X - player.Position.X
			dy := pn.Portals[arrivalID].Pos.Y - player.Position.Y
			player.Position.X = pn.Portals[arrivalID].Pos.X
			player.Position.Y = pn.Portals[arrivalID].Pos.Y
			if player.Hook != nil {
				player.Hook.End.Add(dx, dy)
			}
			return true
		}
	}
	return false
}
