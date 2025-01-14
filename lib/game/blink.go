package game

import (
	"math"
	"time"
)

func (p *Player) BlinkTick() {
	if p.Blinking {
		progress := time.Since(p.BlinkedAt).Seconds() / BlinkDuration
		if progress >= 0.5 && !p.Blinked {
			dx, dy := blinkDistance*math.Cos(p.Angle), blinkDistance*math.Sin(p.Angle)
			p.Position.Add(dx, dy)
			if p.Hook != nil {
				p.Hook.End.Add(dx, dy)
			}
			p.Blinked = true
		}
		if progress >= 1 {
			p.Blinking = false
			p.Blinked = false
		}
	}
}

func (p *Player) HandleBlink() {
	if p.Status != PlayerStatusDead && !p.Blinking && time.Since(p.BlinkedAt).Seconds() >= BlinkCooldown {
		p.Blinking = true
		p.BlinkedAt = time.Now()
	}
}
