package main

import (
	"log/slog"
	"time"

	warsgame "wars/lib/game"
	"wars/lib/messages"
)

func (srv *server) broadcastState() {
	b, err := messages.New(messages.SrvMsgGameState,
		&messages.GameStateMsg{
			Game: &warsgame.Game{Players: srv.game.Players, CId: srv.game.CId},
			Time: time.Now().UnixMilli(),
		}).MarshalMsg(nil)
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
