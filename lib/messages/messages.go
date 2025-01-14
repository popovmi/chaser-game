package messages

import (
	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
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
	ClMsgMove
	ClMsgRotate
	ClMsgBrake
	ClMsgTeleport
	ClMsgBlink
	ClMsgHook
	ClMsgBoost
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

type Empty struct {
}

type YourIDMsg struct {
	ID string `msg:"id"`
}

type JoinGameMsg struct {
	Name string `msg:"name"`
}

type MoveMsg struct {
	Dir string `msg:"dir"`
}

type RotateMsg struct {
	Dir game.Direction `msg:"dir"`
}

type BoostMsg struct {
	Boosting bool `msg:"boosting"`
}

type GameStateMsg struct {
	Game *game.Game `msg:"g"`
}
