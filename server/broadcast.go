package main

import (
	"encoding/binary"
	"log/slog"

	"wars/game"
	"wars/messages"
)

func (srv *server) broadcastState() {

	msg := messages.New(
		messages.SrvMsgGameState,
		&messages.GameStateMsg{
			State: &game.State{
				Players:       srv.game.State.Players,
				PortalNetwork: srv.game.State.PortalNetwork.Short(),
			},
		},
	)

	b, err := msg.MarshalMsg(nil)
	if err != nil {
		slog.Error("could not marshal state", "error", err.Error())
		return
	}

	srv.broadcastUDP(b)
}

func (srv *server) broadcastUDP(b []byte) {
	for _, client := range srv.clients {
		go func() {
			err := client.sendUDPBytes(b)
			if err != nil {
				slog.Error("could not broadcast UDP packet to client", "clientID", client.ID, "error", err.Error())
			}
		}()
	}
}

func (srv *server) broadcastTCP(b []byte) {
	size := uint32(len(b))
	buf := make([]byte, 4+size)
	binary.BigEndian.PutUint32(buf[:4], size)
	copy(buf[4:], b)

	for _, client := range srv.clients {
		go func() {
			err := client.sendTCPBytes(buf)
			if err != nil {
				slog.Error("could not broadcast TCP packet to client", "clientID", client.ID, "error", err.Error())
			}
		}()
	}
}
