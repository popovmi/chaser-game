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

		c.clientID = userIdMsg.ID
		c.screen = screenMain

	case messages.SrvMsgYouJoined:
		state, err := messages.Unmarshal(&game.Game{}, msg.B)
		if err != nil {
			return err
		}
		c.game = state
		for _, player := range c.game.Players {
			c.CreatePlayerImages(player)
		}
		c.drawPortals()
		c.drawBricks()
		c.openUDPConnection()
		c.screen = screenGame

	case messages.SrvMsgPlayerJoined:
		player, err := messages.Unmarshal(&game.Player{}, msg.B)
		if err != nil {
			return err
		}
		c.game.Players[player.ID] = player
		c.CreatePlayerImages(player)

	case messages.SrvMsgGameState:
		state, err := messages.Unmarshal(&messages.GameStateMsg{}, msg.B)
		if err != nil {
			return err
		}

		c.game.ChaserID = state.Game.ChaserID
		for k, player := range c.game.Players {
			if updatedPlayer, ok := state.Game.Players[k]; ok {
				player.JoinedAt = updatedPlayer.JoinedAt
				player.LastChasedAt = updatedPlayer.LastChasedAt
				player.ChaseCount = updatedPlayer.ChaseCount
				player.Position = updatedPlayer.Position
				player.Velocity = updatedPlayer.Velocity
				player.Angle = updatedPlayer.Angle
				player.MoveDir = updatedPlayer.MoveDir
				player.TurnDir = updatedPlayer.TurnDir
				player.Strafing = updatedPlayer.Strafing
				player.Hook = updatedPlayer.Hook
				player.HookedAt = updatedPlayer.HookedAt
				player.IsHooked = updatedPlayer.IsHooked
				player.CaughtByID = updatedPlayer.CaughtByID
				player.Blinking = updatedPlayer.Blinking
				player.BlinkedAt = updatedPlayer.BlinkedAt
				player.Blinked = updatedPlayer.Blinked
				delete(state.Game.Players, k)
			} else {
				delete(c.game.Players, k)
				delete(c.playerImages, k)
			}
		}
		for k, player := range state.Game.Players {
			c.game.Players[k] = player
			c.CreatePlayerImages(player)
		}
		c.game.PreviousTick = time.Now().UnixMilli()
		c.moveCamera()

	case messages.SrvMsgPlayerMoved:
		movedMsg, err := messages.Unmarshal(&messages.PlayerMovedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if movedMsg.ID != c.clientID {
			c.game.Players[movedMsg.ID].HandleMove(movedMsg.Dir)
		}

	case messages.SrvMsgPlayerTurned:
		turnedMsg, err := messages.Unmarshal(&messages.PlayerTurnedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if turnedMsg.ID != c.clientID {
			c.game.Players[turnedMsg.ID].HandleTurn(turnedMsg.Dir)
		}

	case messages.SrvMsgPlayerStrafed:
		strafedMsg, err := messages.Unmarshal(&messages.PlayerStrafedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if strafedMsg.ID != c.clientID {
			c.game.Players[strafedMsg.ID].HandleStrafe(strafedMsg.Strafing)
		}

	case messages.SrvMsgPlayerTeleported:
		portedMsg, err := messages.Unmarshal(&messages.PlayerTeleportedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if portedMsg.ID != c.clientID {
			if c.game.Teleport(portedMsg.ID) {
				c.audio.playPortal()
			}
		}

	case messages.SrvMsgPlayerBlinked:
		blinkedMsg, err := messages.Unmarshal(&messages.PlayerBlinkedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if blinkedMsg.ID != c.clientID {
			c.game.Players[blinkedMsg.ID].HandleBlink()
		}

	case messages.SrvMsgPlayerHooked:
		hookedMsg, err := messages.Unmarshal(&messages.PlayerHookedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if hookedMsg.ID != c.clientID {
			c.game.Players[hookedMsg.ID].UseHook()
		}

	default:
		return nil
	}

	return nil
}
