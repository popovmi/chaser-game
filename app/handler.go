package main

import (
	"log/slog"
	"time"

	"wars/lib/game"
	"wars/lib/messages"
)

func (c *gameClient) handleMessage(msg messages.Message) error {
	switch msg.T {
	case messages.SrvMsgYourID:
		userIdMsg, err := messages.Unmarshal(&messages.YourIDMsg{}, msg.B)
		if err != nil {
			return err
		}

		slog.Info("received user id", "ID", userIdMsg.ID)

		c.id = userIdMsg.ID
		c.ui.screen = screenMain

	case messages.SrvMsgYouJoined:
		state, err := messages.Unmarshal(&warsgame.Game{}, msg.B)
		if err != nil {
			return err
		}
		c.game = state
		for _, player := range c.game.Players {
			c.createPlayerImages(player, player.Color.ToColorRGBA())
		}
		c.drawPortals()
		c.drawBricks()
		c.openUDPConnection()
		c.ui.screen = screenGame

	case messages.SrvMsgPlayerJoined:
		player, err := messages.Unmarshal(&warsgame.Player{}, msg.B)
		if err != nil {
			return err
		}
		c.game.Players[player.ID] = player
		c.createPlayerImages(player, player.Color.ToColorRGBA())

	case messages.SrvMsgGameState:
		state, err := messages.Unmarshal(&messages.GameStateMsg{}, msg.B)
		if err != nil {
			return err
		}
		c.mu.Lock()
		defer c.mu.Unlock()

		c.ping = int(time.Now().UnixMilli() - state.Time)
		c.game.CId = state.Game.CId

		for k, player := range c.game.Players {
			if updatedPlayer, ok := state.Game.Players[k]; ok {
				player.Pos = updatedPlayer.Pos
				player.Vel = updatedPlayer.Vel
				player.Direction = updatedPlayer.Direction
				player.ChaseCount = updatedPlayer.ChaseCount
				player.JoinedAt = updatedPlayer.JoinedAt
				player.LastChasedAt = updatedPlayer.LastChasedAt
				player.BlinkedAt = updatedPlayer.BlinkedAt
				player.Blinking = updatedPlayer.Blinking
				player.IsHooked = updatedPlayer.IsHooked
				player.CaughtByID = updatedPlayer.CaughtByID
				player.Hook = updatedPlayer.Hook
				player.HookedAt = updatedPlayer.HookedAt
				delete(state.Game.Players, k)
			} else {
				delete(c.game.Players, k)
				delete(c.ui.playerImgs, k)
			}
		}
		for k, player := range state.Game.Players {
			c.game.Players[k] = player
			c.createPlayerImages(player, player.Color.ToColorRGBA())
		}
		c.moveCamera()

	case messages.SrvMsgPlayerMoved:
		movedMsg, err := messages.Unmarshal(&messages.PlayerMovedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if movedMsg.ID != c.id {
			c.game.Players[movedMsg.ID].ChangeDirection(movedMsg.Dir)
		}

	case messages.SrvMsgPlayerBraked:
		brakedMsg, err := messages.Unmarshal(&messages.PlayerBrakedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if brakedMsg.ID != c.id {
			c.game.Players[brakedMsg.ID].Brake()
		}

	case messages.SrvMsgPlayerTeleported:
		portedMsg, err := messages.Unmarshal(&messages.PlayerTeleportedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if portedMsg.ID != c.id {
			if c.game.Teleport(portedMsg.ID) {
				c.audio.playPortal()
			}
		}

	case messages.SrvMsgPlayerBlinked:
		blinkedMsg, err := messages.Unmarshal(&messages.PlayerBlinkedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if blinkedMsg.ID != c.id {
			c.game.Blink(blinkedMsg.ID)
		}

	case messages.SrvMsgPlayerHooked:
		hookedMsg, err := messages.Unmarshal(&messages.PlayerHookedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if hookedMsg.ID != c.id {
			c.game.Players[hookedMsg.ID].ThrowHook()
		}

	default:
		return nil
	}

	return nil
}
