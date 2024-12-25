package main

import (
	"image"
	"wars/app/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	warsgame "wars/lib/game"
)

func (c *gameClient) Update() error {
	c.tps = ebiten.ActualTPS()
	switch c.ui.screen {
	case screenMain:
		return c.updateMainScreen()
	case screenGame:
		return c.updateGameScreen()
	default:
		return nil
	}
}

func (c *gameClient) updateMainScreen() error {
	if c.ui.nameInput == nil {
		c.ui.nameInput = components.NewTextField(
			image.Rect(
				c.ui.windowW/2,
				c.ui.windowH/2-warsgame.TextFieldHeight,
				c.ui.windowW/2+warsgame.TextFieldWidth,
				c.ui.windowH/2,
			),
			false,
			FontFace22,
		)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if c.ui.nameInput.Contains(x, y) {
			c.ui.nameInput.Focus()
			c.ui.nameInput.SetSelectionStartByCursorPosition(x, y)
		} else {
			c.ui.nameInput.Blur()
		}
	}

	if err := c.ui.nameInput.Update(); err != nil {
		return err
	}

	x, y := ebiten.CursorPosition()
	if c.ui.nameInput.Contains(x, y) {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	if c.ui.nameInput.IsFocused() && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if err := c.joinGame(c.ui.nameInput.Value()); err != nil {
			return err
		}
	}

	return nil
}

func (c *gameClient) updateGameScreen() error {
	if c.game.Players[c.id] == nil {
		return nil
	}

	err := c.handleMovement()
	if err != nil {
		return err
	}

	wallHits, touches := c.game.Tick()
	for id := range wallHits {
		c.audio.wallHit.play(id)
	}
	for id := range touches {
		c.audio.touch.play(id)
	}

	c.moveCamera()

	for _, p := range c.game.Players {
		if !p.Touchable() {
			ut, ok := c.ui.untouchableTimers[p.ID]
			if !ok {
				c.ui.untouchableTimers[p.ID] = &untouchableTimer{0, true}
				ut = c.ui.untouchableTimers[p.ID]
			}
			ut.t++
			if ut.t > 10 {
				ut.t = 0
				ut.visible = !ut.visible
			}
		} else {
			delete(c.ui.untouchableTimers, p.ID)
		}
	}

	return nil
}
