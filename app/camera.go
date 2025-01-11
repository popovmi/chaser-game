package main

import (
	"time"

	"wars/lib/game"
)

func (c *gameClient) moveCamera() {
	p, ok := c.game.Players[c.clientID]

	if ok && c.windowW < game.FieldWidth+100 {
		if p.Status == game.PlayerStatusDead {
			timeSinceDeath := time.Since(p.DeadAt).Seconds()
			if timeSinceDeath > 0.5 {
				screenCenterX := float64(game.FieldWidth)/2 - float64(c.windowW)/2
				t := (timeSinceDeath - 0.5) / 2.0 // 1.0 - время интерполяции
				if t > 1.0 {
					t = 1.0
				}
				targetX := clamp(p.DeathPos.X-float64(c.windowW)/2, -50, float64(game.FieldWidth-c.windowW)+50)
				c.cameraX = lerp(targetX, screenCenterX, t)
			} else {
				targetX := p.DeathPos.X - float64(c.windowW)/2
				c.cameraX = clamp(targetX, -50, float64(game.FieldWidth-c.windowW)+50)
			}
		} else {
			targetX := p.Position.X - float64(c.windowW)/2
			c.cameraX = clamp(targetX, -50, float64(game.FieldWidth-c.windowW)+50)
		}
	} else {
		c.cameraX = 0
	}

	if ok && c.windowH < game.FieldHeight+100 {
		if p.Status == game.PlayerStatusDead {
			timeSinceDeath := time.Since(p.DeadAt).Seconds()
			if timeSinceDeath > 0.5 {
				screenCenterY := float64(game.FieldHeight)/2 - float64(c.windowH)/2
				t := (timeSinceDeath - 0.5) / 1.0 // 1.0 - время интерполяции
				if t > 1.0 {
					t = 1.0
				}
				targetY := clamp(p.DeathPos.Y-float64(c.windowH)/2, -50, float64(game.FieldHeight-c.windowH)+50)
				c.cameraY = lerp(targetY, screenCenterY, t)
			} else {
				targetY := p.DeathPos.Y - float64(c.windowH)/2
				c.cameraY = clamp(targetY, -50, float64(game.FieldHeight-c.windowH)+50)
			}
		} else {
			targetY := p.Position.Y - float64(c.windowH)/2
			c.cameraY = clamp(targetY, -50, float64(game.FieldHeight-c.windowH)+50)
		}
	} else {
		c.cameraY = 0
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
