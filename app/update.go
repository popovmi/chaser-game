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
	animatePortalIds := make([]string, 0)
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
		c.playerImages[p.ID].animation.Update()
		if p.Teleporting {
			animatePortalIds = append(animatePortalIds, p.DepPortalID, p.ArrPortalID)
		}
		for _, id := range animatePortalIds {
			c.portalAnimations[id].Update()
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

func (c *gameClient) rotate(dir game.RotationDirection) {
	p, ok := c.game.Players[c.clientID]
	if ok && dir != p.RotationDir {
		go func() {
			err := c.sendTCPWithBody(messages.ClMsgRotate, &messages.RotateMsg{Dir: dir})
			if err != nil {
				slog.Error("could not send move", "error", err.Error())
			}
		}()
		p.HandleRotate(dir)
	}
}

func (c *gameClient) move(dir string) {
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

func (c *gameClient) teleport() {
	p, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgTeleport)
			if err != nil {
				slog.Error("could not send teleport", "error", err.Error())
			}
		}()
		if c.game.PortalNetwork.Teleport(p) {
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

func (c *gameClient) brake() {
	p, ok := c.game.Players[c.clientID]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgBrake)
			if err != nil {
				slog.Error("could not send brake", "error", err.Error())
			}
		}()
		p.Brake()
	}
}

func (c *gameClient) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		c.teleport()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.blink()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		c.hook()
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		c.brake()
	}

	c.rotate(getRotateDirection())
	c.move(getMoveDirection())
}

func getRotateDirection() game.RotationDirection {
	var rotateDir game.RotationDirection
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rotateDir = game.RotationNegative
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rotateDir = game.RotationPositive
	}
	return rotateDir
}

func getMoveDirection() string {
	h := chooseDir(isRight(), isLeft(), "r", "l")
	v := chooseDir(isUp(), isDown(), "u", "d")
	return v + h
}

func chooseDir(m, om bool, d, od string) string {
	if m != om {
		if m {
			return d
		}
		if om {
			return od
		}
	}
	return ""
}

func isUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyW)
}

func isDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS)
}

func isLeft() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func isRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyD)
}
