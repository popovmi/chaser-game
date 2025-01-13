package messages

import (
	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
)

//go:generate msgp

type MessageType int

const (
	ClMsgHello MessageType = iota
	ClMsgJoinGame
	ClMsgMove
	ClMsgRotate
	ClMsgBrake
	ClMsgTeleport
	ClMsgBlink
	ClMsgHook
	ClMsgBoost
)

const (
	SrvMsgYourID MessageType = iota
	SrvMsgYouJoined
	SrvMsgPlayerJoined
	SrvMsgGameState
	SrvMsgPlayerMoved
	SrvMsgPlayerRotated
	SrvMsgPlayerBraked
	SrvMsgPlayerTeleported
	SrvMsgPlayerBlinked
	SrvMsgPlayerHooked
	SrvMsgPlayerBoosted
)

type Message struct {
	T MessageType `msg:"type"`
	B msgp.Raw    `msg:"body"`
}

type ClientUDPMessage struct {
	Message
	ID string `msg:"id"`
}

func New(t MessageType, data msgp.Marshaler) *Message {
	body, err := (data).MarshalMsg(nil)
	if err != nil {
		return nil
	}
	return &Message{T: t, B: body}
}

func UDP(t MessageType, id string, data msgp.Marshaler) *ClientUDPMessage {
	m := New(t, data)
	return &ClientUDPMessage{Message: *m, ID: id}
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

type PlayerMovedMsg struct {
	ID  string `msg:"id"`
	Dir string `msg:"dir"`
}

type PlayerRotatedMsg struct {
	ID  string         `msg:"id"`
	Dir game.Direction `msg:"dir"`
}

type PlayerTeleportedMsg struct {
	ID string `msg:"id"`
}

type PlayerBlinkedMsg struct {
	ID string `msg:"id"`
}

type PlayerHookedMsg struct {
	ID string `msg:"id"`
}

type PlayerBrakedMsg struct {
	ID string `msg:"id"`
}

type PlayerBoostedMsg struct {
	ID       string `msg:"id"`
	Boosting bool   `msg:"boosting"`
}

type GameStateMsg struct {
	Game *game.Game `msg:"g"`
}
