package main

import (
	"log/slog"
	"time"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

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
	case messages.ClMsgMove,
		messages.ClMsgRotate,
		messages.ClMsgBoost,
		messages.ClMsgTeleport,
		messages.ClMsgBlink,
		messages.ClMsgHook,
		messages.ClMsgBrake:
		{
			srvMsg = &messages.ClientMessage{}
		}
	default:
	}

	if srvMsg == nil {
		return nil
	}
	_, err := srvMsg.UnmarshalMsg(msg.B)
	if err != nil {
		slog.Info("could not unmarshal srvMsg", "err", err.Error())
		return err
	}

	if clientMessage, ok := srvMsg.(*messages.ClientMessage); ok {
		c.handleClientMessage(clientMessage)
		return nil
	}

	switch msg.T {
	case messages.SrvMsgYourID:
		c.handleYourId(srvMsg)

	case messages.SrvMsgYouJoined:
		c.handleYouJoined(srvMsg.(*game.Game))

	case messages.SrvMsgPlayerJoined:
		c.handlePlayerJoined(srvMsg.(*game.Player))

	case messages.SrvMsgGameState:
		c.handleGameState(srvMsg.(*messages.GameStateMsg))

	default:
		return nil
	}
	return nil
}

func (c *gameClient) handleYourId(srvMsg msgp.Unmarshaler) {
	msg := srvMsg.(*messages.YourIDMsg)
	slog.Info("received user id", "ID", msg.ID)
	c.clientID = msg.ID
	c.screen = screenMain
}

func (c *gameClient) handleYouJoined(msg *game.Game) {
	c.game = msg
	for _, player := range c.game.Players {
		c.createPlayerImages(player)
	}
	c.drawBricks()
	c.createPortalsAnimations()
	c.openUDPConnection()
	c.screen = screenGame
}

func (c *gameClient) handlePlayerJoined(msg *game.Player) {
	c.createPlayerImages(msg)
	c.game.Players[msg.ID] = msg
}

func (c *gameClient) handleGameState(msg *messages.GameStateMsg) {
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
	c.game.PreviousTick = time.Now().UnixMilli()
	c.moveCamera()
}

func (c *gameClient) handleClientMessage(msg *messages.ClientMessage) {
	clientId := msg.ID
	if clientId == c.clientID {
		return
	}
	player := c.game.Players[clientId]
	if player == nil {
		return
	}

	var clMsg msgp.Unmarshaler
	if msg.T == messages.ClMsgMove {
		clMsg = &messages.MoveMsg{}
	}
	if msg.T == messages.ClMsgBoost {
		clMsg = &messages.BoostMsg{}
	}
	if msg.T == messages.ClMsgRotate {
		clMsg = &messages.RotateMsg{}
	}
	if clMsg != nil {
		_, err := clMsg.UnmarshalMsg(msg.B)
		if err != nil {
			slog.Info("could not unmarshal clMsg", "err", err.Error())
			return
		}
	}

	switch msg.T {
	case messages.ClMsgMove:
		player.HandleMove(clMsg.(*messages.MoveMsg).Dir)

	case messages.ClMsgRotate:
		player.HandleRotate(clMsg.(*messages.RotateMsg).Dir)

	case messages.ClMsgBoost:
		player.HandleBoost(clMsg.(*messages.BoostMsg).Boosting)

	case messages.ClMsgTeleport:
		if c.game.PortalNetwork.Teleport(player) {
			c.audio.playPortal()
		}

	case messages.ClMsgBlink:
		player.HandleBlink()

	case messages.ClMsgHook:
		player.UseHook()

	case messages.ClMsgBrake:
		player.Brake()

	default:
	}
}
