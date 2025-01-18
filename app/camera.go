package main

import (
	"time"

	"wars/game"
)

func (c *gameClient) translateCamera() {
	p, ok := c.game.State.Players[c.clientID]

	if !ok {
		c.cameraX = 0
		c.cameraY = 0
		return
	}

	halfScreenWidth := float64(c.windowW) / 2
	halfFieldWidth := float64(game.FieldWidth) / 2
	screenCenterX := halfFieldWidth - halfScreenWidth
	minX := float64(-50)
	maxX := float64(game.FieldWidth-c.windowW) + 50

	if c.windowW < game.FieldWidth+100 {
		if p.Status == game.PlayerStatusPreparing {
			c.cameraX = screenCenterX
		} else if p.Status == game.PlayerStatusDead {
			deathPosSubX := p.DeathPosition.X - halfFieldWidth
			timeSinceDeath := time.Since(*p.DeadAt).Seconds()
			if timeSinceDeath <= 0.5 {
				c.cameraX = clamp(deathPosSubX, minX, maxX)
			} else {
				t := (timeSinceDeath - 0.5) / 2.0
				if t > 1.0 {
					t = 1.0
				}
				targetX := clamp(deathPosSubX, minX, maxX)
				c.cameraX = lerp(targetX, screenCenterX, t)
			}
		} else {
			targetX := p.Position.X - halfScreenWidth
			targetX = clamp(targetX, minX, maxX)
			//c.cameraX = lerp(c.cameraX, targetX, factor) // как правильно
			c.cameraX = targetX
		}
	} else {
		c.cameraX = 0
	}

	halfScreenHeight := float64(c.windowH) / 2
	halfFieldHeight := float64(game.FieldHeight) / 2
	screenCenterY := halfFieldHeight - halfScreenHeight
	minY := float64(-50)
	maxY := float64(game.FieldHeight-c.windowH) + 50

	if c.windowH < game.FieldHeight+100 {
		if p.Status == game.PlayerStatusPreparing {
			c.cameraY = screenCenterY
		} else if p.Status == game.PlayerStatusDead {
			deathPosY := p.DeathPosition.Y - halfScreenHeight
			timeSinceDeath := time.Since(*p.DeadAt).Seconds()
			if timeSinceDeath <= 0.5 {
				c.cameraY = clamp(deathPosY, minY, maxY)
			} else {
				t := (timeSinceDeath - 0.5) / 2.0
				if t > 1.0 {
					t = 1.0
				}
				targetY := clamp(deathPosY, minY, maxY)
				c.cameraY = lerp(targetY, screenCenterY, t)
			}
		} else {
			targetY := p.Position.Y - halfScreenHeight
			targetY = clamp(targetY, minY, maxY)
			//c.cameraY = lerp(c.cameraY, targetY, factor) // как правильно
			c.cameraY = targetY
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
