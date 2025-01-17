package main

import (
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"wars/app/components"
	"wars/game"
	"wars/messages"
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
	state := c.game.State
	_, ok := state.Players[c.clientID]
	if !ok {
		return nil
	}

	animatePortalIds := make([]string, 0)
	for _, p := range state.Players {
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
		//if animaion, ok :=
		c.playerImages[p.ID].animation.Update()
		if p.Teleporting {
			animatePortalIds = append(animatePortalIds, p.FromPortalID, p.ToPortalID)
		}
		for _, id := range animatePortalIds {
			c.portalAnimations[id].Update()
		}
	}

	c.HandleInput()
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

func (c *gameClient) HandleInput() {
	p, ok := c.game.State.Players[c.clientID]
	if !ok {
		return
	}
	var commands game.Commands
	if !p.Teleporting && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		commands = append(commands, game.Command{Action: game.CommandActionTeleport, PlayerID: p.ID})
	}
	if !p.Blinking && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		commands = append(commands, game.Command{Action: game.CommandActionBlink, PlayerID: p.ID})
	}
	if p.Hook == nil && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		commands = append(commands, game.Command{Action: game.CommandActionHook, PlayerID: p.ID})
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		commands = append(commands, game.Command{Action: game.CommandActionBrake, PlayerID: p.ID})
	}
	if dir := getMoveDirection(); dir != p.MoveDirection {
		commands = append(commands, game.Command{Action: game.CommandActionMove, PlayerID: p.ID, Payload: dir})
	}
	if dir := getRotateDirection(); dir != p.RotationDirection {
		commands = append(commands, game.Command{Action: game.CommandActionRotate, PlayerID: p.ID, Payload: dir})
	}
	if boosting := ebiten.IsKeyPressed(ebiten.KeyUp); boosting != p.Boosting {
		commands = append(commands, game.Command{Action: game.CommandActionBoost, PlayerID: p.ID, Payload: boosting})
	}
	if len(commands) > 0 {
		c.game.AddCommands(commands)
		err := c.sendUDPWithBody(messages.ClMsgInGameCommandPack, commands)
		if err != nil {
			slog.Error("could not send command", "error", err.Error())
		}
	}

}

func getRotateDirection() game.Direction {
	var rotateDir game.Direction
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rotateDir = game.DirectionNegative
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rotateDir = game.DirectionPositive
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
