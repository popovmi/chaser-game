package main

import (
	"log/slog"
	"time"

	"github.com/tinylib/msgp/msgp"

	"wars/game"
	"wars/messages"
)

func (c *gameClient) handleMessage(msg *messages.Message, expired bool) error {
	var srvMsg msgp.Unmarshaler
	switch msg.T {
	case messages.SrvMsgPong:
		c.ping = time.Since(c.lastPingTime)
		return nil
	case messages.SrvMsgYourID:
		srvMsg = &messages.YourIDMsg{}
	case messages.SrvMsgYouJoined:
		srvMsg = &game.Game{}
	case messages.SrvMsgPlayerJoined:
		srvMsg = &game.Player{}
	case messages.SrvMsgGameState:
		if expired {
			return nil
		}
		srvMsg = &messages.GameStateMsg{}
	case messages.ClMsgInGameCommand:
		command := game.Command{}
		_, err := command.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
		if command.PlayerID != c.clientID {
			c.game.AddCommand(command)
		}
		return nil
	case messages.ClMsgInGameCommandPack:
		commands := game.Commands{}
		_, err := commands.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
		slog.Info("got commands: ", commands)
		for _, command := range commands {
			if command.PlayerID != c.clientID {
				c.game.AddCommand(command)
			}
		}
		return nil
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
	c.game.State = msg.State
	slog.Info("got state", "state", c.game.State)
	for _, player := range c.game.State.Players {
		c.createPlayerImages(player)
	}
	c.drawBricks()
	c.createPortalsAnimations()
	c.openUDPConnection()
	c.game.Start()
	err := c.sendUDPWithBody(
		messages.ClMsgInGameCommand,
		game.Command{Action: game.CommandActionReady, PlayerID: c.clientID},
	)
	if err != nil {
		slog.Info("could not send ready command", "err", err.Error())
	}
	c.screen = screenGame
	go c.startPingRoutine()
}

func (c *gameClient) handlePlayerJoined(msg *game.Player) {
	c.createPlayerImages(msg)
	c.game.State.Players[msg.ID] = msg
}

func (c *gameClient) handleGameState(msg *messages.GameStateMsg) {
	gState := c.game.State
	for k, player := range gState.Players {
		if updatedPlayer, ok := msg.State.Players[k]; ok {
			player.Set(updatedPlayer)
			delete(msg.State.Players, k)
		} else {
			delete(gState.Players, k)
			delete(c.playerImages, k)
		}
	}
	for k, player := range msg.State.Players {
		gState.Players[k] = player
		c.createPlayerImages(player)
	}
	for k, link := range msg.State.PortalNetwork.Links {
		gState.PortalNetwork.Links[k].LastUsedMap = link.LastUsedMap
	}
	for k, portal := range msg.State.PortalNetwork.Portals {
		gState.PortalNetwork.Portals[k].LastUsedAt = portal.LastUsedAt
	}
	c.game.LastTick = time.Now()
}
