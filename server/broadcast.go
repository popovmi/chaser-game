package main

import (
	"log/slog"

	"wars/lib/game"
	"wars/lib/messages"
)

func (srv *server) broadcastState() {
	b, err := messages.New(
		messages.SrvMsgGameState,
		&messages.GameStateMsg{
			Game: &game.Game{
				Players: srv.game.Players,
				PortalNetwork: &game.PortalNetwork{
					Links: srv.game.PortalNetwork.Links,
				},
			},
		},
	).MarshalMsg(nil)

	if err != nil {
		slog.Error("could not marshal state", err.Error())
		return
	}
	srv.broadcastUDP(b)
}

func (srv *server) broadcastUDP(b []byte) {
	for _, player := range srv.game.Players {
		go func() { _ = srv.clients[player.ID].sendUDPBytes(b) }()
	}
}
