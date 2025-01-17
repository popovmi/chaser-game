package messages

import (
	"github.com/tinylib/msgp/msgp"

	"wars/game"
)

//go:generate msgp

type MessageType int

const (
	SrvMsgPong MessageType = iota
	SrvMsgYourID
	SrvMsgYouJoined
	SrvMsgPlayerJoined
	SrvMsgGameState

	ClMsgHello
	ClMsgPing
	ClMsgJoinGame
	ClMsgInGameCommand
	ClMsgInGameCommandPack
)

type Message struct {
	T MessageType `msg:"type"`
	B msgp.Raw    `msg:"body"`
}

type ClientMessage struct {
	ID string `msg:"id"`
	Message
}

func New(t MessageType, data msgp.Marshaler) *Message {
	body, err := (data).MarshalMsg(nil)
	if err != nil {
		return nil
	}
	return &Message{T: t, B: body}
}

func UDP(t MessageType, id string, data msgp.Marshaler) *ClientMessage {
	m := New(t, data)
	return &ClientMessage{Message: *m, ID: id}
}

type YourIDMsg struct {
	ID string `msg:"id"`
}

type JoinGameMsg struct {
	Name string `msg:"name"`
}

type GameStateMsg struct {
	State *game.State `msg:"s"`
}

type Empty struct{}
