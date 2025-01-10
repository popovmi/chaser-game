package game

import (
	"math"
	"time"

	"chaser/lib/vector"
)

//go:generate msgp

type Hook struct {
	End             vector.Vector2D `msg:"end"`
	Vel             vector.Vector2D `msg:"vel"`
	CurrentDistance float64         `msg:"current_distance,omitempty"`
	Stuck           bool            `msg:"stuck,omitempty"`
	IsReturning     bool            `msg:"is_returning,omitempty"`
	CaughtPlayerID  string          `msg:"caught_player_id,omitempty"`
}

func (p *Player) HookTick(dt float64, players map[string]*Player) {
	if p.Hook == nil {
		return
	}

	if p.Hook.Stuck {
		vel := vector.NewVector2D(p.Hook.End.X, p.Hook.End.Y)
		vel.SubV(p.Position)
		vel.Normalize()
		vel.Mul(hookBackwardVelocity)
		p.Velocity = vel
		p.Step(dt)

		if p.IsHookDone() {
			return
		}
	} else {
		target, targetExists := players[p.Hook.CaughtPlayerID]

		if p.Hook.IsReturning {
			vel := vector.NewVector2D(p.Position.X, p.Position.Y)
			vel.SubV(p.Hook.End)
			vel.Normalize()
			vel.Mul(hookBackwardVelocity * dt)
			p.Hook.Vel = vel
			p.Hook.End.AddV(p.Hook.Vel)
		} else {
			p.Hook.End.Add(p.Hook.Vel.X*dt, p.Hook.Vel.Y*dt)
			p.Hook.CurrentDistance += hookVelocity * dt
			if p.Hook.CurrentDistance >= hookDistance {
				p.Hook.IsReturning = true
			}
		}

		if !targetExists {
			p.HookPlayer(players)
			target, targetExists = players[p.Hook.CaughtPlayerID]
		}

		if targetExists {
			target.Position.X = p.Hook.End.X
			target.Position.Y = p.Hook.End.Y
		} else if !p.Hook.IsReturning {
			p.Hook.Clamp()
		}

		if p.IsHookDone() {
			if targetExists {
				target.IsHooked = false
				target.CaughtByID = ""
			}
		}
	}
}

func (p *Player) UseHook() {
	if p.Hook != nil && !p.Hook.IsReturning && p.Hook.CurrentDistance >= hookMinDistance {
		p.Hook.IsReturning = true
		return
	}
	if p.Hook == nil && time.Since(p.HookedAt).Seconds() >= HookCooldown {
		p.Hook = &Hook{
			End: vector.NewVector2D(p.Position.X, p.Position.Y),
			Vel: vector.NewVector2D(math.Cos(p.Angle)*hookVelocity, math.Sin(p.Angle)*hookVelocity),
		}
	}
}

func (p *Player) IsHookDone() bool {
	if !p.Hook.Stuck && !p.Hook.IsReturning {
		return false
	}
	d := p.Hook.End.Distance(p.Position)
	if d < Radius || (p.Hook.CaughtPlayerID != "" && d < 2*Radius) {
		p.Hook = nil
		p.HookedAt = time.Now()
		return true
	}
	return false
}

func (h *Hook) Clamp() {
	if h.End.X < 0 {
		h.End.X = 0
		h.Stuck = true
	}
	if h.End.X > FieldWidth {
		h.End.X = FieldWidth
		h.Stuck = true
	}
	if h.End.Y < 0 {
		h.End.Y = 0
		h.Stuck = true
	}
	if h.End.Y > FieldHeight {
		h.End.Y = FieldHeight
		h.Stuck = true
	}
	if h.Stuck {
		h.IsReturning = true
	}
}

func (p *Player) HookPlayer(players map[string]*Player) bool {
	for _, target := range players {
		if target.ID == p.ID {
			continue
		}

		distance := target.Position.Distance(p.Hook.End)

		if distance < Radius {
			p.Hook.CaughtPlayerID = target.ID
			p.Hook.IsReturning = true
			target.IsHooked = true
			target.CaughtByID = p.ID
			target.Velocity.X = 0
			target.Velocity.Y = 0
			target.MoveDir = ""
			return true
		}
	}

	return false
}
