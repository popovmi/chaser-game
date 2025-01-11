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
			c.createPlayerImages(player)
		}
		c.drawBricks()
		c.createPortalsAnimations()
		c.openUDPConnection()
		c.screen = screenGame

	case messages.SrvMsgPlayerJoined:
		player, err := messages.Unmarshal(&game.Player{}, msg.B)
		if err != nil {
			return err
		}
		c.game.Players[player.ID] = player
		c.createPlayerImages(player)

	case messages.SrvMsgGameState:
		state, err := messages.Unmarshal(&messages.GameStateMsg{}, msg.B)
		if err != nil {
			return err
		}
		for k, player := range c.game.Players {
			if updatedPlayer, ok := state.Game.Players[k]; ok {
				*player = *updatedPlayer
				delete(state.Game.Players, k)
			} else {
				delete(c.game.Players, k)
				delete(c.playerImages, k)
			}
		}
		for k, player := range state.Game.Players {
			c.game.Players[k] = player
			c.createPlayerImages(player)
		}
		for k, link := range state.Game.PortalNetwork.Links {
			c.game.PortalNetwork.Links[k].LastUsed = link.LastUsed
		}
		for k, portal := range state.Game.PortalNetwork.Portals {
			*c.game.PortalNetwork.Portals[k] = *portal
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

	case messages.SrvMsgPlayerRotated:
		rotatedMsg, err := messages.Unmarshal(&messages.PlayerRotatedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if rotatedMsg.ID != c.clientID {
			c.game.Players[rotatedMsg.ID].HandleRotate(rotatedMsg.Dir)
		}

	case messages.SrvMsgPlayerTeleported:
		portedMsg, err := messages.Unmarshal(&messages.PlayerTeleportedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if portedMsg.ID != c.clientID {
			if c.game.PortalNetwork.Teleport(c.game.Players[portedMsg.ID]) {
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

	case messages.SrvMsgPlayerBraked:
		brakedMsg, err := messages.Unmarshal(&messages.PlayerBrakedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if brakedMsg.ID != c.clientID {
			c.game.Players[brakedMsg.ID].Brake()
		}

	default:
		return nil
	}

	return nil
}
