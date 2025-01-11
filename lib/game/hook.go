package game

import (
	"math"
	"time"

	"wars/lib/vector"
)

//go:generate msgp

type Hook struct {
	End            vector.Vector2D `msg:"end"`
	Vel            vector.Vector2D `msg:"vel"`
	Stuck          bool            `msg:"stuck"`
	IsReturning    bool            `msg:"is_returning"`
	CaughtPlayerID string          `msg:"caught_player_id"`
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
		return
	}

	p.RotateHook()
	p.Hook.End.Add(p.Hook.Vel.X*dt, p.Hook.Vel.Y*dt)
	if !p.Hook.IsReturning && p.Hook.End.Distance(p.Position) >= MaxHookLength {
		p.Hook.IsReturning = true
	}

	target, targetExists := players[p.Hook.CaughtPlayerID]
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

func (p *Player) UseHook() {
	if p.Hook != nil && !p.Hook.IsReturning {
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
}

func (p *Player) HookPlayer(players map[string]*Player) bool {
	for _, target := range players {
		if target.ID == p.ID || !target.Touchable() {
			continue
		}

		distance := target.Position.Distance(p.Hook.End)
		if distance < Radius {
			target.takeHookHit(p)
			return true
		}
	}

	return false
}

func (p *Player) takeHookHit(hookedPlayer *Player) {
	p.HP -= hookDamage
	if p.HP <= 0 {
		p.die()
		hookedPlayer.Kills += 1
		hookedPlayer.Hook.CaughtPlayerID = ""
		hookedPlayer.Hook.IsReturning = true
		return
	}

	hookedPlayer.Hook.CaughtPlayerID = p.ID
	hookedPlayer.Hook.IsReturning = true
	p.IsHooked = true
	p.CaughtByID = hookedPlayer.ID
	p.Velocity.X = 0
	p.Velocity.Y = 0
	p.MoveDir = ""
}

func (p *Player) RotateHook() {
	if p.Hook == nil || p.Hook.Stuck {
		return
	}
	length := p.HookLength()

	p.Hook.End.X = p.Position.X + length*math.Cos(p.Angle)
	p.Hook.End.Y = p.Position.Y + length*math.Sin(p.Angle)

	if p.Hook.IsReturning {
		p.Hook.Vel.X = -math.Cos(p.Angle) * hookBackwardVelocity
		p.Hook.Vel.Y = -math.Sin(p.Angle) * hookBackwardVelocity
	} else {
		p.Hook.Vel.X = math.Cos(p.Angle) * hookVelocity
		p.Hook.Vel.Y = math.Sin(p.Angle) * hookVelocity
	}
}

func (p *Player) HookLength() float64 {
	v := vector.NewVector2D(p.Hook.End.X, p.Hook.End.Y)
	v.SubV(p.Position)
	return v.Length()
}
