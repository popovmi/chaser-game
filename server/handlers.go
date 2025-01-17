package main

import (
	"errors"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/game"
	"wars/messages"
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

	c := &srvClient{Player: game.NewPlayer(), ip: host, tcp: conn}
	srv.clients[c.ID] = c

	defer func() {
		c.tcp.Close()
		delete(srv.clients, c.ID)
		srv.game.DeletePlayer(c.ID)
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
	switch msg.T {
	case messages.ClMsgPing:
		return c.handlePing()

	case messages.ClMsgHello:
		slog.Info("new UDP client", "ID", c.ID, "IP", c.udpAddr.IP.String())
		return nil

	case messages.ClMsgJoinGame:
		clMsg := &messages.JoinGameMsg{}
		_, err := clMsg.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
		return srv.join(c, clMsg)

	case messages.ClMsgInGameCommandPack:
		commands := game.Commands{}
		_, err := commands.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
		slog.Debug("got command pack", "playerID", c.ID, "commands", commands)
		srv.game.AddCommands(commands)
		broadcastMsg := messages.New(msg.T, commands)
		b, err := broadcastMsg.MarshalMsg(nil)
		if err != nil {
			slog.Error("could not encode message", "error", err.Error(), "msg", broadcastMsg)
		}
		srv.broadcastUDP(b)

	case messages.ClMsgInGameCommand:
		command := game.Command{}
		_, err := command.UnmarshalMsg(msg.B)
		if err != nil {
			return err
		}
		srv.game.AddCommand(command)

		broadcastMsg := messages.New(msg.T, command)
		b, err := broadcastMsg.MarshalMsg(nil)
		if err != nil {
			slog.Error("could not encode message", "error", err.Error(), "msg", broadcastMsg)
		}
		srv.broadcastUDP(b)

	default:
		return nil
	}

	return nil
}

func (c *srvClient) handlePing() error {
	err := c.sendTCP(messages.SrvMsgPong)
	if err != nil {
		slog.Error("could not send Pong", "error", err.Error())
		return err
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
	srv.game.Join(c.Player)

	if err := c.sendTCPWithBody(messages.SrvMsgYouJoined, srv.game); err != nil {
		return err
	}
	for _, p := range srv.clients {
		if p.ID != c.ID {
			if err := srv.clients[p.ID].sendTCPWithBody(messages.SrvMsgPlayerJoined, c.Player); err != nil {
				slog.Info("could not send joined player", "joinedID", c.ID, "targetID", p.ID)
			}
		}
	}
	return nil
}
