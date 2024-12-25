package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"wars/lib/messages"
)

func (c *gameClient) joinGame(name string) error {
	c.ui.screen = screenWait
	err := c.sendTCPWithBody(messages.ClMsgJoinGame, &messages.JoinGameMsg{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *gameClient) move(dir string) error {
	p, ok := c.game.Players[c.id]
	if ok && dir != p.Direction {
		go func() {
			err := c.sendTCPWithBody(messages.ClMsgMove, &messages.MoveMsg{Dir: dir})
			if err != nil {
				slog.Error("could not send move", "error", err.Error())
			}
		}()
		p.ChangeDirection(dir)
	}
	return nil
}

func (c *gameClient) brake() error {
	p, ok := c.game.Players[c.id]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgBrake)
			if err != nil {
				slog.Error("could not send brake", "error", err.Error())
			}
		}()
		p.Brake()
	}
	return nil
}

func (c *gameClient) teleport() error {
	_, ok := c.game.Players[c.id]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgTeleport)
			if err != nil {
				slog.Error("could not send teleport", "error", err.Error())
			}
		}()
		if c.game.Teleport(c.id) {
			c.audio.playPortal()
		}
	}
	return nil
}

func (c *gameClient) blink() error {
	_, ok := c.game.Players[c.id]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgBlink)
			if err != nil {
				slog.Error("could not send teleport", "error", err.Error())
			}
		}()
		c.game.Blink(c.id)
	}
	return nil
}

func (c *gameClient) hook() error {
	p, ok := c.game.Players[c.id]
	if ok {
		go func() {
			err := c.sendTCP(messages.ClMsgHook)
			if err != nil {
				slog.Error("could not send hook", "error", err.Error())
			}
		}()
		p.ThrowHook()
	}
	return nil
}

func (c *gameClient) handleMovement() error {
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		return c.teleport()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return c.blink()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return c.hook()
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		return c.brake()
	}

	h := chooseDir(isRight(), isLeft(), "r", "l")
	v := chooseDir(isUp(), isDown(), "u", "d")
	return c.move(h + v)
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
	return ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
}

func isDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
}

func isLeft() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
}

func isRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
}
