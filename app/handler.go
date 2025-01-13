package main

import (
	"log/slog"
	"time"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

type handler func(msg msgp.Unmarshaler)

func (c *gameClient) handleMessage(msg *messages.Message) error {
	var srvMsg msgp.Unmarshaler
	switch msg.T {
	case messages.SrvMsgYourID:
		srvMsg = &messages.YourIDMsg{}
	case messages.SrvMsgYouJoined:
		srvMsg = &game.Game{}
	case messages.SrvMsgPlayerJoined:
		srvMsg = &game.Player{}
	case messages.SrvMsgGameState:
		srvMsg = &messages.GameStateMsg{}
	case messages.SrvMsgPlayerMoved:
		srvMsg = &messages.PlayerMovedMsg{}
	case messages.SrvMsgPlayerRotated:
		srvMsg = &messages.PlayerRotatedMsg{}
	case messages.SrvMsgPlayerTeleported:
		srvMsg = &messages.PlayerTeleportedMsg{}
	case messages.SrvMsgPlayerBlinked:
		srvMsg = &messages.PlayerBlinkedMsg{}
	case messages.SrvMsgPlayerHooked:
		srvMsg = &messages.PlayerHookedMsg{}
	case messages.SrvMsgPlayerBraked:
		srvMsg = &messages.PlayerBrakedMsg{}
	case messages.SrvMsgPlayerBoosted:
		srvMsg = &messages.PlayerBoostedMsg{}
	}

	if srvMsg != nil {
		_, err := srvMsg.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
	}

	map[messages.MessageType]handler{
		messages.SrvMsgYourID:           c.handleYourId,
		messages.SrvMsgYouJoined:        c.handleYouJoined,
		messages.SrvMsgPlayerJoined:     c.handlePlayerJoined,
		messages.SrvMsgGameState:        c.handleGameState,
		messages.SrvMsgPlayerMoved:      c.handlePlayerMoved,
		messages.SrvMsgPlayerRotated:    c.handlePlayerRotated,
		messages.SrvMsgPlayerTeleported: c.handlePlayerTeleported,
		messages.SrvMsgPlayerBlinked:    c.handlePlayerBlinked,
		messages.SrvMsgPlayerHooked:     c.handlePlayerHooked,
		messages.SrvMsgPlayerBraked:     c.handlePlayerBraked,
		messages.SrvMsgPlayerBoosted:    c.handlePlayerBoosted,
	}[msg.T](srvMsg)

	return nil
}

func (c *gameClient) handleYourId(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.YourIDMsg)
	slog.Info("received user id", "ID", msg.ID)
	c.clientID = msg.ID
	c.screen = screenMain
}

func (c *gameClient) handleYouJoined(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*game.Game)
	c.game = msg
	for _, player := range c.game.Players {
		c.createPlayerImages(player)
	}
	c.drawBricks()
	c.createPortalsAnimations()
	c.openUDPConnection()
	c.screen = screenGame
}

func (c *gameClient) handlePlayerJoined(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*game.Player)
	c.createPlayerImages(msg)
	c.game.Players[msg.ID] = msg
}

func (c *gameClient) handleGameState(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.GameStateMsg)
	if c.game == nil {
		return
	}
	for k, player := range c.game.Players {
		if updatedPlayer, ok := msg.Game.Players[k]; ok {
			*player = *updatedPlayer
			delete(msg.Game.Players, k)
		} else {
			delete(c.game.Players, k)
			delete(c.playerImages, k)
		}
	}
	for k, player := range msg.Game.Players {
		c.game.Players[k] = player
		c.createPlayerImages(player)
	}
	for k, link := range msg.Game.PortalNetwork.Links {
		c.game.PortalNetwork.Links[k].LastUsed = link.LastUsed
	}
	for k, portal := range msg.Game.PortalNetwork.Portals {
		*c.game.PortalNetwork.Portals[k] = *portal
	}
	c.game.PreviousTick = time.Now().UnixMilli()
	c.moveCamera()
}

func (c *gameClient) handlePlayerMoved(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerMovedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].HandleMove(msg.Dir)
	}
}

func (c *gameClient) handlePlayerRotated(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerRotatedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].HandleRotate(msg.Dir)
	}
}

func (c *gameClient) handlePlayerTeleported(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerTeleportedMsg)
	if msg.ID != c.clientID {
		if c.game.PortalNetwork.Teleport(c.game.Players[msg.ID]) {
			c.audio.playPortal()
		}
	}
}

func (c *gameClient) handlePlayerBlinked(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerBlinkedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].HandleBlink()
	}

}

func (c *gameClient) handlePlayerHooked(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerHookedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].UseHook()
	}
}

func (c *gameClient) handlePlayerBraked(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerBrakedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].Brake()
	}
}

func (c *gameClient) handlePlayerBoosted(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.PlayerBoostedMsg)
	if msg.ID != c.clientID {
		c.game.Players[msg.ID].HandleBoost(msg.Boosting)
	}
}
