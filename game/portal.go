package game

import (
	"log/slog"
	"sync"
	"time"
)

//go:generate msgp

type Portal struct {
	ID         string     `msg:"id" json:"ID"`
	LinkID     string     `msg:"lid,omitempty" json:"LinkID,omitempty"`
	LastUsedAt *time.Time `msg:"lua,omitempty" json:"LastUsedAt,omitempty"`
	Position   *Vector    `msg:"pos,omitempty" json:"Position,omitempty"`
}

func newPortal(id, linkID string, x, y float64) *Portal {
	return &Portal{ID: id, LinkID: linkID, Position: NewVector(x, y)}
}

type PortalLink struct {
	ID          string                `msg:"id" json:"ID"`
	PortalIDs   []string              `msg:"pls,omitempty" json:"PortalIDs,omitempty"`
	LastUsedMap map[string]*time.Time `msg:"lum,omitempty" json:"LastUsedMap,omitempty"`
}

func newPortalLink(id string) *PortalLink {
	return &PortalLink{ID: id, LastUsedMap: make(map[string]*time.Time)}
}

type PortalNetwork struct {
	Portals map[string]*Portal     `msg:"pls,omitempty" json:"Portals,omitempty"`
	Links   map[string]*PortalLink `msg:"lks,omitempty" json:"Links,omitempty"`

	mu sync.Mutex
}

func newPortalNetwork(portals map[string]*Portal, links map[string]*PortalLink) *PortalNetwork {
	return &PortalNetwork{Portals: portals, Links: links}
}

func (pn *PortalNetwork) CanUsePortal(player *Player) (bool, *Portal, *time.Duration) {
	if player.Status == PlayerStatusPreparing || player.HookedBy != "" {
		return false, nil, nil
	}
	var portal *Portal
	for _, plr := range pn.Portals {
		if player.touchingPortal(plr) {
			portal = plr
			break
		}
	}
	if portal == nil {
		return false, nil, nil
	}
	link := pn.Links[portal.LinkID]
	lastUsed, ok := link.LastUsedMap[player.ID]
	if !ok {
		return true, portal, nil
	}
	cooldown := time.Since(*lastUsed)
	if cooldown.Seconds() < PortalCooldown {
		return false, portal, &cooldown
	}
	return true, portal, nil
}

func (pn *PortalNetwork) teleport(player *Player) bool {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	can, fromPortal, _ := pn.CanUsePortal(player)
	if !can {
		return false
	}
	link := pn.Links[fromPortal.LinkID]
	for _, toPortalID := range link.PortalIDs {
		if toPortalID != fromPortal.ID {
			now := time.Now()
			player.Teleporting = true
			player.FromPortalID = fromPortal.ID
			player.ToPortalID = toPortalID
			player.TeleportedAt = &now
			fromPortal.LastUsedAt = &now
			pn.Portals[toPortalID].LastUsedAt = &now
			link.LastUsedMap[player.ID] = &now
			return true
		}
	}
	return false
}

func (pn *PortalNetwork) tick(player *Player) {
	if !player.Teleporting {
		return
	}
	fromPort := pn.Portals[player.FromPortalID]
	toPort := pn.Portals[player.ToPortalID]
	slog.Debug("teleporting portal", "fromPort", fromPort, "toPort", toPort)
	link := pn.Links[fromPort.LinkID]
	progress := time.Since(*link.LastUsedMap[player.ID]).Seconds() / TeleportDuration
	if progress >= 0.5 && !player.Teleported {
		dx := toPort.Position.X - player.Position.X
		dy := toPort.Position.Y - player.Position.Y
		player.translatePosition(dx, dy)
		player.Teleported = true
	}
	if progress >= 1 {
		player.Teleporting = false
		player.Teleported = false
		player.FromPortalID = ""
		player.ToPortalID = ""
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
			ID:          link.ID,
			LastUsedMap: link.LastUsedMap,
		}
	}
	return short
}
