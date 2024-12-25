package main

import warsgame "wars/lib/game"

func (c *gameClient) moveCamera() {
	if c.ui.windowW < warsgame.FieldWidth+100 {
		c.ui.cameraX = c.game.Players[c.id].Pos.X - float64(c.ui.windowW)/2
		c.ui.cameraX = clamp(c.ui.cameraX, -50, float64(warsgame.FieldWidth-c.ui.windowW)+50)
	} else {
		c.ui.cameraX = 0
	}

	if c.ui.windowH < warsgame.FieldHeight+100 {
		c.ui.cameraY = c.game.Players[c.id].Pos.Y - float64(c.ui.windowH)/2
		c.ui.cameraY = clamp(c.ui.cameraY, -50, float64(warsgame.FieldHeight-c.ui.windowH)+50)
	} else {
		c.ui.cameraY = 0
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
