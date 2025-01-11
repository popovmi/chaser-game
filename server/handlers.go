package main

import (
	"errors"
	"log/slog"
	"net"
	"strconv"

	"github.com/tinylib/msgp/msgp"
	
	"wars/lib/game"
	"wars/lib/messages"
)

func (srv *server) handleUDPData(addr *net.UDPAddr, data []byte, n int) error {
	msg, err := messages.Unmarshal(&messages.ClientUDPMessage{}, data[:n])
	if err != nil {
		slog.Error("could not decode UDP data", err)
		return err
	}

	ip := addr.IP.String()
	id := msg.ID
	c, ok := srv.clients[id]
	if !ok {
		slog.Error("client not found", "ID", id, "udpIP", ip)
		return errors.New("client not found")
	}
	if c.ip != ip {
		slog.Error("client ip not match", "ID", id, "tcpIP", c.ip, "udpIP", ip)
		return errors.New("client ip not match")
	}

	c.udpAddr = addr
	c.udp = srv.udp

	return srv.handleMessage(c, msg.Message)
}

func (srv *server) handleTCPConnection(conn net.Conn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		slog.Error("could not split host and port", err.Error())
		return
	}

	id := srv.game.Counter.Add(1)
	c := &srvClient{Player: game.NewPlayer(strconv.FormatUint(id, 10)), ip: host, tcp: conn}
	srv.clients[c.ID] = c

	defer func() {
		c.tcp.Close()
		delete(srv.clients, c.ID)
		srv.game.RemovePlayer(c.ID)
		srv.colors[c.Color] = false
	}()

	err = c.sendTCPWithBody(messages.SrvMsgYourID, &messages.YourIDMsg{ID: c.ID})
	if err != nil {
		return
	}

	for {
		var msg messages.Message
		if err := msgp.Decode(c.tcp, &msg); err != nil {
			slog.Error("could not decode TCP message", err.Error())
			return
		}

		err := srv.handleMessage(c, msg)
		if err != nil {
			return
		}
	}
}

func (srv *server) handleMessage(c *srvClient, msg messages.Message) error {
	switch msg.T {

	case messages.ClMsgHello:
		slog.Info("new UDP client", "ID", c.ID, "IP", c.udpAddr.IP.String())

	case messages.ClMsgJoinGame:
		joinReq, err := messages.Unmarshal(&messages.JoinGameMsg{}, msg.B)
		if err != nil {
			return err
		}

		slog.Info(
			"new join request",
			"ID", c.ID,
			"IP", c.ip,
			"name", joinReq.Name,
		)

		c.Name = joinReq.Name
		for clrKey, picked := range srv.colors {
			if !picked {
				srv.colors[clrKey] = true
				c.Color = clrKey
				break
			}
		}
		srv.game.AddPlayer(c.Player)
		if err = c.sendTCPWithBody(messages.SrvMsgYouJoined, srv.game); err != nil {
			return err
		}
		for _, p := range srv.game.Players {
			if p.ID != c.ID {
				if err = srv.clients[p.ID].sendTCPWithBody(messages.SrvMsgPlayerJoined, c.Player); err != nil {
					return err
				}
			}
		}

	case messages.ClMsgMove:
		moveReq, err := messages.Unmarshal(&messages.MoveMsg{}, msg.B)
		if err != nil {
			return err
		}
		c.HandleMove(moveReq.Dir)
		movedMsg := &messages.PlayerMovedMsg{ID: c.ID, Dir: moveReq.Dir}
		b, err := messages.New(messages.SrvMsgPlayerMoved, movedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	case messages.ClMsgRotate:
		rotateReq, err := messages.Unmarshal(&messages.RotateMsg{}, msg.B)
		if err != nil {
			return err
		}
		c.HandleTurn(rotateReq.Dir)
		turnedMsg := &messages.PlayerTurnedMsg{ID: c.ID, Dir: rotateReq.Dir}
		b, err := messages.New(messages.SrvMsgPlayerMoved, turnedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	case messages.ClMsgTeleport:
		srv.game.Teleport(c.ID)
		portedMsg := &messages.PlayerTeleportedMsg{ID: c.ID}
		b, err := messages.New(messages.SrvMsgPlayerTeleported, portedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	case messages.ClMsgBlink:
		c.HandleBlink()
		blinkedMsg := &messages.PlayerBlinkedMsg{ID: c.ID}
		b, err := messages.New(messages.SrvMsgPlayerBlinked, blinkedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	case messages.ClMsgHook:
		c.UseHook()
		hookedMsg := &messages.PlayerHookedMsg{ID: c.ID}
		b, err := messages.New(messages.SrvMsgPlayerHooked, hookedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	case messages.ClMsgBrake:
		c.Brake()
		brakedMsg := &messages.PlayerBrakedMsg{ID: c.ID}
		b, err := messages.New(messages.SrvMsgPlayerHooked, brakedMsg).MarshalMsg(nil)
		if err != nil {
			slog.Error("could not marshal updates", err)
			return err
		}
		srv.broadcastUDP(b)

	default:
		return nil
	}

	return nil
}
