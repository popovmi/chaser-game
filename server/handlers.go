package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

func (srv *server) handleUDPData(addr *net.UDPAddr, data []byte, n int) error {
	msg := &messages.ClientMessage{}
	_, err := msg.UnmarshalMsg(data[:n])
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

	num := srv.game.Counter.Add(1)
	id := fmt.Sprintf("p%d", num)
	c := &srvClient{Player: game.NewPlayer(id), ip: host, tcp: conn}
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
			slog.Error("could not decode TCP message", "error", err.Error())
			return
		}

		err := srv.handleMessage(c, msg)
		if err != nil {
			slog.Error("could not handle TCP message", "error", err.Error())
			return
		}
	}
}

func (srv *server) handleMessage(c *srvClient, msg messages.Message) error {
	var clMsg msgp.Unmarshaler
	needsUnmarshal := false
	needsBroadcast := true
	switch msg.T {
	case messages.ClMsgHello:
		slog.Info("new UDP client", "ID", c.ID, "IP", c.udpAddr.IP.String())
		return nil

	case messages.ClMsgJoinGame:
		clMsg = &messages.JoinGameMsg{}
		needsUnmarshal = true
		needsBroadcast = false

	case messages.ClMsgMove:
		clMsg = &messages.MoveMsg{}
		needsUnmarshal = true

	case messages.ClMsgRotate:
		clMsg = &messages.RotateMsg{}
		needsUnmarshal = true

	case messages.ClMsgBoost:
		clMsg = &messages.BoostMsg{}
		needsUnmarshal = true

	case messages.ClMsgTeleport, messages.ClMsgBlink, messages.ClMsgHook, messages.ClMsgBrake:
		break

	default:
		return nil
	}

	if needsBroadcast {
		go func() {
			broadcastMsg := messages.New(msg.T, &messages.ClientMessage{ID: c.ID, Message: msg})
			b, err := broadcastMsg.MarshalMsg(nil)
			if err != nil {
				slog.Error("could not encode message", "error", err.Error(), "msg", broadcastMsg)
			}
			srv.broadcastUDP(b)
		}()
	}

	if needsUnmarshal {
		_, err := clMsg.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
	}

	switch msg.T {
	case messages.ClMsgHello:
		slog.Info("new UDP client", "ID", c.ID, "IP", c.udpAddr.IP.String())
	case messages.ClMsgJoinGame:
		return srv.join(c, clMsg.(*messages.JoinGameMsg))
	case messages.ClMsgMove:
		return srv.move(c, clMsg.(*messages.MoveMsg))
	case messages.ClMsgRotate:
		return srv.rotate(c, clMsg.(*messages.RotateMsg))
	case messages.ClMsgTeleport:
		return srv.teleport(c)
	case messages.ClMsgBlink:
		return srv.blink(c)
	case messages.ClMsgHook:
		return srv.hook(c)
	case messages.ClMsgBoost:
		return srv.boost(c, clMsg.(*messages.BoostMsg))
	case messages.ClMsgBrake:
		return srv.brake(c)
	default:
		return nil
	}

	return nil
}

func (srv *server) join(c *srvClient, msg *messages.JoinGameMsg) error {
	slog.Info(
		"new join request",
		"ID", c.ID,
		"IP", c.ip,
		"name", msg.Name,
	)
	c.Name = msg.Name
	for clrKey, picked := range srv.colors {
		if !picked {
			srv.colors[clrKey] = true
			c.Color = clrKey
			break
		}
	}
	srv.game.AddPlayer(c.Player)
	if err := c.sendTCPWithBody(messages.SrvMsgYouJoined, srv.game); err != nil {
		return err
	}
	for _, p := range srv.game.Players {
		if p.ID != c.ID {
			go func() {
				if err := srv.clients[p.ID].sendTCPWithBody(messages.SrvMsgPlayerJoined, c.Player); err != nil {
					slog.Info("could not send joined player", "joinedID", c.ID, "targetID", p.ID)
				}
			}()
		}
	}
	return nil
}

func (srv *server) move(c *srvClient, msg *messages.MoveMsg) error {
	c.HandleMove(msg.Dir)
	return nil
}

func (srv *server) rotate(c *srvClient, msg *messages.RotateMsg) error {
	c.HandleRotate(msg.Dir)
	return nil
}

func (srv *server) teleport(c *srvClient) error {
	srv.game.PortalNetwork.Teleport(srv.game.Players[c.ID])
	return nil
}

func (srv *server) blink(c *srvClient) error {
	c.HandleBlink()
	return nil
}

func (srv *server) hook(c *srvClient) error {
	c.UseHook()
	return nil
}

func (srv *server) brake(c *srvClient) error {
	c.Brake()
	return nil
}

func (srv *server) boost(c *srvClient, msg *messages.BoostMsg) error {
	c.HandleBoost(msg.Boosting)
	return nil
}
