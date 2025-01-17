package game

import (
	"math"
	"time"
)

//go:generate msgp

type Hook struct {
	EndPosition    *Vector `msg:"e,omitempty" json:"EndPosition,omitempty"`
	Velocity       *Vector `msg:"v,omitempty" json:"Velocity,omitempty"`
	Stuck          bool    `msg:"s,omitempty" json:"Stuck,omitempty"`
	Returning      bool    `msg:"r,omitempty" json:"Returning,omitempty"`
	HookedPlayerID string  `msg:"hpi,omitempty" json:"HookedPlayerID,omitempty"`
}

func (p *Player) HookLength() float64 {
	return p.Hook.EndPosition.DistanceTo(p.Position)
}

func (p *Player) takeHookHit(hookedPlayer *Player) {
	p.HP -= hookDamage
	if p.HP <= 0 {
		p.die()
		hookedPlayer.Kills += 1
		hookedPlayer.Hook.HookedPlayerID = ""
		hookedPlayer.Hook.Returning = true
		return
	}

	hookedPlayer.Hook.HookedPlayerID = p.ID
	hookedPlayer.Hook.Returning = true
	p.HookedBy = hookedPlayer.ID
	p.Velocity.X = 0
	p.Velocity.Y = 0
	p.MoveDirection = ""
	p.Boosting = false
}

func (p *Player) rotateHook() {
	if p.Hook == nil || p.Hook.Stuck {
		return
	}
	length := p.HookLength()

	p.Hook.EndPosition = p.Position.Clone()
	p.Hook.EndPosition.Translate(length*math.Cos(p.Angle), length*math.Sin(p.Angle))

	vel := float64(hookVelocity)
	if p.Hook.Returning {
		vel = -hookBackwardVelocity
	}
	p.Hook.Velocity = NewVector(math.Cos(p.Angle)*vel, math.Sin(p.Angle)*vel)
}

func (p *Player) hookPlayer(players map[string]*Player) bool {
	for _, target := range players {
		if target.ID == p.ID || !target.Touchable() {
			continue
		}

		distance := target.Position.DistanceTo(p.Hook.EndPosition)
		if distance < Radius {
			target.takeHookHit(p)
			return true
		}
	}

	return false
}

func (p *Player) isHookDone() bool {
	if !p.Hook.Stuck && !p.Hook.Returning {
		return false
	}
	d := p.Hook.EndPosition.DistanceTo(p.Position)
	if d < Radius || (p.Hook.HookedPlayerID != "" && d < 2*Radius) {
		t := time.Now()
		p.Hook = nil
		p.UsedHookAt = &t
		return true
	}
	return false
}

func (p *Player) hookClamp() {
	h := p.Hook
	if h.EndPosition.X < 0 {
		h.EndPosition.X = 0
		h.Stuck = true
	}
	if h.EndPosition.X > FieldWidth {
		h.EndPosition.X = FieldWidth
		h.Stuck = true
	}
	if h.EndPosition.Y < 0 {
		h.EndPosition.Y = 0
		h.Stuck = true
	}
	if h.EndPosition.Y > FieldHeight {
		h.EndPosition.Y = FieldHeight
		h.Stuck = true
	}
	if h.Stuck {
		h.Velocity = NewVector(0, 0)
	}
}

func (p *Player) hookTick(dt float64, players map[string]*Player) {
	if p.Hook == nil {
		return
	}

	if p.Hook.Stuck {
		vel := p.Hook.EndPosition.Clone()
		vel.Subtract(p.Position)
		vel.Normalize()
		vel.Scale(hookBackwardVelocity)
		p.Velocity = vel
		p.step(dt)

		if p.isHookDone() {
			return
		}
		return
	}

	p.rotateHook()
	p.Hook.EndPosition.Translate(p.Hook.Velocity.X*dt, p.Hook.Velocity.Y*dt)
	if !p.Hook.Returning && p.Hook.EndPosition.DistanceTo(p.Position) >= MaxHookLength {
		p.Hook.Returning = true
	}

	if p.Hook.HookedPlayerID == "" {
		p.hookPlayer(players)
	}

	target, targetExists := players[p.Hook.HookedPlayerID]

	if targetExists {
		target.Position = p.Hook.EndPosition.Clone()
	} else if !p.Hook.Returning {
		p.hookClamp()
	}

	if p.isHookDone() {
		if targetExists {
			target.HookedBy = ""
		}
	}

}
