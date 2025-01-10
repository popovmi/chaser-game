package main

import (
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"wars/app/components"
	"wars/lib/game"
	"wars/lib/messages"
)

func (c *gameClient) Update() error {
	c.tps = ebiten.ActualTPS()
	switch c.screen {
	case screenMain:
		return c.updateMainScreen()
	case screenGame:
		return c.updateGameScreen()
	default:
		return nil
	}
}

func (c *gameClient) updateMainScreen() error {
	if c.nameInput == nil {
		c.nameInput = components.NewTextField(
			image.Rect(
				c.windowW/2,
				c.windowH/2-30,
				c.windowW/2+282,
				c.windowH/2,
			),
			false,
			FontFace22,
			20, 30,
		)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if c.nameInput.Contains(x, y) {
			c.nameInput.Focus()
			c.nameInput.SetSelectionStartByCursorPosition(x, y)
		} else {
			c.nameInput.Blur()
		}
	}

	if err := c.nameInput.Update(); err != nil {
		return err
	}

	x, y := ebiten.CursorPosition()
	if c.nameInput.Contains(x, y) {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	if c.nameInput.IsFocused() && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if err := c.joinGame(c.nameInput.Value()); err != nil {
			return err
		}
	}

	return nil
}

func (c *gameClient) updateGameScreen() error {
	if c.game.Players[c.clientID] == nil {
		return nil
	}
	c.HandleInput()
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
			ut, ok := c.untouchableTimers[p.ID]
			if !ok {
				c.untouchableTimers[p.ID] = &untouchableTimer{0, true}
				ut = c.untouchableTimers[p.ID]
			}
			ut.t++
			if ut.t > 10 {
				ut.t = 0
				ut.visible = !ut.visible
			}
		} else {
			delete(c.untouchableTimers, p.ID)
		}
	}
	return nil
}

func (c *gameClient) joinGame(name string) error {
	c.screen = screenWait
	err := c.sendTCPWithBody(messages.ClMsgJoinGame, &messages.JoinGameMsg{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *gameClient) turn(dir game.Direction) {
	p, ok := c.game.Players[c.clientID]
	if ok && dir != p.TurnDir {
		go func() {
			err := c.sendTCPWithBody(messages.ClMsgTurn, &messages.TurnMsg{Dir: dir})
			if err != nil {
				slog.Error("could not send move", "error", err.Error())
			}
		}()
		p.HandleTurn(dir)
	}
}

func (c *gameClient) move(dir game.Direction) {
	p, ok := c.game.Players[c.clientID]
	if ok && dir != p.MoveDir {
		go func() {
			err := c.sendTCPWithBody(messages.ClMsgMove, &messages.MoveMsg{Dir: dir})
			if err != nil {
				slog.Error("could not send move", "error", err.Error())
			}
		}()
		p.HandleMove(dir)
	}
}

func (c *gameClient) strafe(strafing bool) {
	p, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCPWithBody(messages.ClMsgStrafe, &messages.StrafeMsg{Strafing: strafing})
			if err != nil {
				slog.Error("could not send brake", "error", err.Error())
			}
		}()
		p.HandleStrafe(strafing)
	}
}

func (c *gameClient) teleport() {
	_, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgTeleport)
			if err != nil {
				slog.Error("could not send teleport", "error", err.Error())
			}
		}()
		if c.game.Teleport(c.clientID) {
			c.audio.playPortal()
		}
	}
}

func (c *gameClient) blink() {
	p, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgBlink)
			if err != nil {
				slog.Error("could not send teleport", "error", err.Error())
			}
		}()
		p.HandleBlink()
	}
}

func (c *gameClient) hook() {
	p, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgHook)
			if err != nil {
				slog.Error("could not send hook", "error", err.Error())
			}
		}()
		p.UseHook()
	}
}

func (c *gameClient) HandleInput() {
	p := c.game.Players[c.clientID]
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.teleport()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.blink()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		c.hook()
	}

	strafe := false
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		strafe = true
	}
	if p.Strafing != strafe {
		c.strafe(strafe)
	}

	moveDir := game.ZeroDir
	if p.Strafing {
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			moveDir = game.PositiveDir
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			moveDir = game.NegativeDir
		}
		c.move(moveDir)
	} else {
		turnDir := game.ZeroDir
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			turnDir = game.PositiveDir
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			turnDir = game.NegativeDir
		}
		c.turn(turnDir)

		if ebiten.IsKeyPressed(ebiten.KeyW) {
			moveDir = game.PositiveDir
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			moveDir = game.NegativeDir
		}
		c.move(moveDir)
	}
}
