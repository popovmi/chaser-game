package main

import (
	"encoding/binary"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

type srvClient struct {
	*game.Player

	ip      string
	tcp     net.Conn
	udpAddr *net.UDPAddr
	udp     *net.UDPConn
}

func (c *srvClient) sendTCP(t messages.MessageType) error {
	msg := messages.Message{T: t}
	data, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}
	size := uint32(len(data))
	buf := make([]byte, 4+size)
	binary.BigEndian.PutUint32(buf[:4], size)
	copy(buf[4:], data)
	return c.sendTCPBytes(buf)
}

func (c *srvClient) sendTCPWithBody(t messages.MessageType, body msgp.Marshaler) error {
	msg := messages.New(t, body)
	data, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}

	size := uint32(len(data))
	buf := make([]byte, 4+size)
	binary.BigEndian.PutUint32(buf[:4], size)
	copy(buf[4:], data)

	return c.sendTCPBytes(buf)
}

func (c *srvClient) sendTCPBytes(b []byte) error {
	if _, err := c.tcp.Write(b); err != nil {
		slog.Error("could not send TCP message", err)
		return err
	}
	return nil
}

func (c *srvClient) sendUDPBytes(b []byte) error {
	if c.udp != nil {
		if _, err := c.udp.WriteToUDP(b, c.udpAddr); err != nil {
			slog.Error("could not send UDP message", err)
			return err
		}
	}
	return nil
}
