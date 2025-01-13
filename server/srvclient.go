package main

import (
	"encoding/binary"
	"log/slog"
	"net"
	"sync/atomic"

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

	udpMsgCount atomic.Uint64
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
		const maxUDPSize = 1012

		messageID := c.udpMsgCount.Add(1)
		totalPackets := uint16((len(b) + maxUDPSize - 1) / maxUDPSize)
		for i := uint16(0); i < totalPackets; i++ {
			header := make([]byte, 12)
			binary.BigEndian.PutUint64(header[0:8], messageID)
			binary.BigEndian.PutUint16(header[8:10], totalPackets)
			binary.BigEndian.PutUint16(header[10:12], i)

			start := int(i) * maxUDPSize
			end := start + maxUDPSize
			if end > len(b) {
				end = len(b)
			}
			payload := b[start:end]

			packet := append(header, payload...)
			_, err := c.udp.WriteToUDP(packet, c.udpAddr)
			if err != nil {
				slog.Error("could not send udp packet", "error", err.Error())
			}
		}
		return nil
	}
	return nil
}
