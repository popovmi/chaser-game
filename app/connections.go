package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/messages"
)

func (c *gameClient) openTCPConnection() {
	conn, err := net.Dial("tcp", c.tcpAddr)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	c.TCPConn = conn
	go c.handleTCP()
}

func (c *gameClient) handleTCP() {
	for {

		sizeBuf := make([]byte, 5)
		_, err := io.ReadFull(c.TCPConn, sizeBuf)
		if err != nil {
			slog.Error("could not read TCP message header", "error", err.Error())
			if errors.Is(err, io.EOF) {
				slog.Info("TCP connection closed")
				break
			}
			continue
		}

		size := binary.BigEndian.Uint32(sizeBuf)
		data := make([]byte, size)
		_, err = io.ReadFull(c.TCPConn, data)
		if err != nil {
			slog.Error("could not read TCP message body", "error", err.Error())
			if errors.Is(err, io.EOF) {
				slog.Info("TCP connection closed")
				break
			}
			continue
		}

		msg := &messages.Message{}
		_, err = msg.UnmarshalMsg(data)
		if err != nil {
			slog.Error("could not decode TCP message", "error", err.Error())
			continue
		}

		err = c.handleMessage(msg)
		if err != nil {
			slog.Error("could not handle TCP message", "error", err.Error())
			continue
		}
	}
}

func (c *gameClient) openUDPConnection() {
	udpAddr, err := net.ResolveUDPAddr("udp", c.udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	c.UDPConn = conn
	if err := c.sendUDP(messages.ClMsgHello); err != nil {
		log.Fatal(err)
	}

	go c.handleUDP()
}

func (c *gameClient) handleUDP() {
	for {
		buf := make([]byte, 1024*10)
		n, err := c.UDPConn.Read(buf)
		if err != nil {
			slog.Error("could not read UDP message", "error", err.Error())
			continue
		}
		msg := &messages.Message{}
		_, err = msg.UnmarshalMsg(buf[:n])
		if err != nil {
			slog.Error("could not decode UDP message", "error", err.Error())
			continue
		}
		if err := c.handleMessage(msg); err != nil {
			slog.Error("could not handle UDP message", "error", err.Error())
			continue
		}
	}
}

func (c *gameClient) sendTCP(t messages.MessageType) error {
	return c.sendTCPWithBody(t, &messages.Empty{})
}

func (c *gameClient) sendUDP(t messages.MessageType) error {
	return c.sendUDPWithBody(t, &messages.Empty{})
}

func (c *gameClient) sendTCPWithBody(t messages.MessageType, data msgp.Marshaler) error {
	if err := msgp.Encode(c.TCPConn, messages.New(t, data)); err != nil {
		slog.Error("could not send TCP message", err.Error())
		return err
	}
	return nil
}

func (c *gameClient) sendUDPWithBody(t messages.MessageType, data msgp.Marshaler) error {
	if err := msgp.Encode(c.UDPConn, messages.UDP(t, c.clientID, data)); err != nil {
		slog.Error("could not send UDP message", err.Error())
		return err
	}
	return nil
}

func (c *gameClient) sendMsg(conType string, t messages.MessageType) error {
	switch conType {
	case "tcp":
		return c.sendTCP(t)
	case "udp":
		return c.sendUDP(t)
	default:
		return errors.New("unknown con type")
	}
}

func (c *gameClient) sendMsgWithBody(conType string, t messages.MessageType, data msgp.Marshaler) error {
	switch conType {
	case "tcp":
		return c.sendTCPWithBody(t, data)
	case "udp":
		return c.sendUDPWithBody(t, data)
	default:
		return errors.New("unknown con type")
	}
}
