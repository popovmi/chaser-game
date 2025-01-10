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
	ClMsgTurn
	ClMsgStrafe
	ClMsgTeleport
	ClMsgBlink
	ClMsgHook
)

const (
	SrvMsgYourID MessageType = iota
	SrvMsgYouJoined
	SrvMsgPlayerJoined
	SrvMsgGameState
	SrvMsgPlayerMoved
	SrvMsgPlayerTurned
	SrvMsgPlayerStrafed
	SrvMsgPlayerTeleported
	SrvMsgPlayerBlinked
	SrvMsgPlayerHooked
)

type Message struct {
	T MessageType `msg:"type"`
	B MessageBody `msg:"body"`
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
	Dir game.Direction `msg:"dir"`
}

type TurnMsg struct {
	Dir game.Direction `msg:"dir"`
}

type StrafeMsg struct {
	Strafing bool `msg:"str"`
}

type UdpMoveMsg struct {
	ClientUDPMessage
	MoveMsg
}

type PlayerMovedMsg struct {
	ID  string         `msg:"id"`
	Dir game.Direction `msg:"dir"`
}

type PlayerTurnedMsg struct {
	ID  string         `msg:"id"`
	Dir game.Direction `msg:"dir"`
}

type PlayerStrafedMsg struct {
	ID       string `msg:"id"`
	Strafing bool   `msg:"str"`
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

type GameStateMsg struct {
	Game *game.Game `msg:"g"`
}

func Unmarshal[T msgp.Unmarshaler](msg T, b []byte) (T, error) {
	_, err := msg.UnmarshalMsg(b)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
