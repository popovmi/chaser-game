package main

import (
	"log/slog"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

func (c *gameClient) handleMessage(msg *messages.Message) error {
	var (
		srvMsg, clMsg   msgp.Unmarshaler
		isClientMessage bool
		clientId        string
	)

	switch msg.T {
	case messages.SrvMsgYourID:
		srvMsg = &messages.YourIDMsg{}
	case messages.SrvMsgYouJoined:
		srvMsg = &game.Game{}
	case messages.SrvMsgPlayerJoined:
		srvMsg = &game.Player{}
	case messages.SrvMsgGameState:
		srvMsg = &messages.GameStateMsg{}
	case messages.ClMsgMove:
		isClientMessage = true
		clMsg = &messages.MoveMsg{}
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgRotate:
		isClientMessage = true
		clMsg = &messages.RotateMsg{}
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgBoost:
		isClientMessage = true
		clMsg = &messages.BoostMsg{}
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgTeleport:
		isClientMessage = true
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgBlink:
		isClientMessage = true
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgHook:
		isClientMessage = true
		srvMsg = &messages.ClientMessage{}
	case messages.ClMsgBrake:
		isClientMessage = true
		srvMsg = &messages.ClientMessage{}
	default:
	}

	if srvMsg != nil {
		_, err := srvMsg.UnmarshalMsg(msg.B)
		if err != nil {
			slog.Info("could not unmarshal srvMsg", "msg", msg, "err", err.Error())
			return err
		}
	}

	if isClientMessage {
		clientId = srvMsg.(*messages.ClientMessage).ID
	}

	if clMsg != nil {
		_, err := clMsg.UnmarshalMsg(srvMsg.(*messages.ClientMessage).B)
		if err != nil {
			slog.Info("could not unmarshal clMsg", "srvMsg", srvMsg, "err", err.Error())
			return err
		}
	}

	switch msg.T {
	case messages.SrvMsgYourID:
		c.handleYourId(srvMsg)

	case messages.SrvMsgYouJoined:
		c.handleYouJoined(srvMsg)

	case messages.SrvMsgPlayerJoined:
		c.handlePlayerJoined(srvMsg)

	case messages.SrvMsgGameState:
		c.handleGameState(srvMsg)

	case messages.ClMsgMove:
		c.handlePlayerMoved(clientId, clMsg.(*messages.MoveMsg))

	case messages.ClMsgRotate:
		c.handlePlayerRotated(clientId, clMsg.(*messages.RotateMsg))

	case messages.ClMsgBoost:
		c.handlePlayerBoosted(clientId, clMsg.(*messages.BoostMsg))

	case messages.ClMsgTeleport:
		c.handlePlayerTeleported(clientId)

	case messages.ClMsgBlink:
		c.handlePlayerBlinked(clientId)

	case messages.ClMsgHook:
		c.handlePlayerHooked(clientId)

	case messages.ClMsgBrake:
		c.handlePlayerBraked(clientId)
	}

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
		c.game.PortalNetwork.Portals[k].LastUsedAt = portal.LastUsedAt
	}
	c.game.Tick()
	c.moveCamera()
}

func (c *gameClient) handlePlayerMoved(id string, msg *messages.MoveMsg) {
	if id != c.clientID {
		c.game.Players[id].HandleMove(msg.Dir)
	}
}

func (c *gameClient) handlePlayerRotated(id string, msg *messages.RotateMsg) {
	if id != c.clientID {
		c.game.Players[id].HandleRotate(msg.Dir)
	}
}

func (c *gameClient) handlePlayerTeleported(id string) {
	if id != c.clientID {
		if c.game.PortalNetwork.Teleport(c.game.Players[id]) {
			c.audio.playPortal()
		}
	}
}

func (c *gameClient) handlePlayerBlinked(id string) {
	if id != c.clientID {
		c.game.Players[id].HandleBlink()
	}

}

func (c *gameClient) handlePlayerHooked(id string) {
	if id != c.clientID {
		c.game.Players[id].UseHook()
	}
}

func (c *gameClient) handlePlayerBraked(id string) {
	if id != c.clientID {
		c.game.Players[id].Brake()
	}
}

func (c *gameClient) handlePlayerBoosted(id string, msg *messages.BoostMsg) {
	if id != c.clientID {
		c.game.Players[id].HandleBoost(msg.Boosting)
	}
}
