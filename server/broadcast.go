package main

import (
	"encoding/binary"
	"log/slog"

	"wars/lib/game"
	"wars/lib/messages"
)

func (srv *server) broadcastState() {
	msg := messages.New(
		messages.SrvMsgGameState,
		&messages.GameStateMsg{
			Game: &game.Game{
				Players:       srv.game.Players,
				PortalNetwork: srv.game.PortalNetwork,
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
	for _, player := range srv.game.Players {
		go func() {
			err := srv.clients[player.ID].sendUDPBytes(b)
			if err != nil {
				slog.Error("could not broadcast UDP packet to client", "clientID", player.ID, "error", err.Error())
			}
		}()
	}
}

func (srv *server) broadcastTCP(b []byte) {
	size := uint32(len(b))
	buf := make([]byte, 4+size)
	binary.BigEndian.PutUint32(buf[:4], size)
	copy(buf[4:], b)

	for _, player := range srv.game.Players {
		go func() {
			err := srv.clients[player.ID].sendTCPBytes(buf)
			if err != nil {
				slog.Error("could not broadcast TCP packet to client", "clientID", player.ID, "error", err.Error())
			}
		}()
	}
}
