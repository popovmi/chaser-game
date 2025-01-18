package game

import (
	"log/slog"
	"math"
	"time"
)

func (g *Game) Join(player *Player) {
	player.mu.Lock()
	t := time.Now()
	player.Status = PlayerStatusPreparing
	player.JoinedAt = &t
	player.mu.Unlock()
	g.mu.Lock()
	g.State.Players[player.ID] = player
	g.mu.Unlock()
}

func (g *Game) handlePlayerMove(player *Player, dir string) {
	player.mu.Lock()
	player.MoveDirection = dir
	player.mu.Unlock()
	g.events <- Event{EventActionMoved, player}
}

func (g *Game) handlePlayerRotate(player *Player, dir Direction) {
	player.mu.Lock()
	player.RotationDirection = dir
	player.mu.Unlock()
	g.events <- Event{EventActionRotated, player}
}

func (g *Game) handlePlayerBoost(player *Player, boosting bool) {
	if player.HookedBy == "" && (player.Hook == nil || !player.Hook.Stuck) {
		player.mu.Lock()
		player.Boosting = boosting
		player.mu.Unlock()
		g.events <- Event{EventActionBoosted, player}
	}
}

func (g *Game) handlePlayerBrake(player *Player) {
	player.mu.Lock()
	player.MoveDirection = ""
	player.Velocity.Scale(Braking)
	player.mu.Unlock()
	g.events <- Event{EventActionBraked, player.ID}

}

func (g *Game) handlePlayerTeleport(player *Player) {
	player.mu.Lock()
	defer player.mu.Unlock()
	can, fromPortal, _ := g.State.PortalNetwork.CanUsePortal(player)
	if !can {
		return
	}
	slog.Debug("publish event")
	g.events <- Event{EventActionPortalUsed, player}
	ported := g.State.PortalNetwork.teleport(player, fromPortal)
	if ported {
		g.events <- Event{EventActionTeleported, player.ID}
	}
}

func (g *Game) handlePlayerBlink(player *Player) {
	if !player.Blinking && (player.BlinkedAt == nil || time.Since(*player.BlinkedAt).Seconds() >= BlinkCooldown) {
		player.mu.Lock()
		t := time.Now()
		player.Blinking = true
		player.BlinkedAt = &t
		player.mu.Unlock()
		g.events <- Event{EventActionBlinked, player.ID}
	}
}

func (g *Game) handlePlayerUseHook(player *Player) {
	player.mu.Lock()
	defer player.mu.Unlock()
	if player.Hook != nil && !player.Hook.Returning {
		player.Hook.Returning = true
		g.events <- Event{EventActionReturnedHook, player.ID}
		return
	}
	if player.Hook == nil && (player.UsedHookAt == nil || time.Since(*player.UsedHookAt).Seconds() >= HookCooldown) {
		player.Hook = &Hook{
			EndPosition: player.Position.Clone(),
			Velocity:    NewVector(math.Cos(player.Angle)*hookVelocity, math.Sin(player.Angle)*hookVelocity),
		}
		g.events <- Event{EventActionHooked, player.ID}
	}
}
