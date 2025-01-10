package main

import (
	"chaser/lib/game"
)

func (c *gameClient) moveCamera() {
	p, ok := c.game.Players[c.clientID]

	if ok && c.windowW < game.FieldWidth+100 {
		targetX := p.Position.X - float64(c.windowW)/2
		c.cameraX = clamp(targetX, -50, float64(game.FieldWidth-c.windowW)+50)
	} else {
		c.cameraX = 0
	}

	if ok && c.windowH < game.FieldHeight+100 {
		targetY := p.Position.Y - float64(c.windowH)/2
		c.cameraY = clamp(targetY, -50, float64(game.FieldHeight-c.windowH)+50)
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
