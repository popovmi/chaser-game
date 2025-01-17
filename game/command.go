package game

import (
	"fmt"
	"log/slog"
	"time"
)

//go:generate msgp

type CommandAction byte

const (
	CommandActionNone CommandAction = iota
	CommandActionReady
	CommandActionMove
	CommandActionRotate
	CommandActionBrake
	CommandActionTeleport
	CommandActionBlink
	CommandActionHook
	CommandActionBoost
)

func (ca CommandAction) String() string {
	switch ca {
	case CommandActionNone:
		return "None"
	case CommandActionReady:
		return "Ready"
	case CommandActionMove:
		return "Move"
	case CommandActionRotate:
		return "Rotate"
	case CommandActionBrake:
		return "Brake"
	case CommandActionTeleport:
		return "Teleport"
	case CommandActionBlink:
		return "Blink"
	case CommandActionHook:
		return "Hook"
	case CommandActionBoost:
		return "Boost"
	default:
		return fmt.Sprintf("Unknown(%d)", ca)
	}
}

func (ca CommandAction) MarshalJSON() (data []byte, err error) {
	return []byte(fmt.Sprintf("\"%s\"", ca.String())), nil
}

type Commands []Command

type Command struct {
	Action   CommandAction `msg:"act,omitempty" json:"Action,omitempty"`
	PlayerID string        `msg:"pid,omitempty" json:"PlayerID,omitempty"`
	Payload  interface{}   `msg:"pld,omitempty" json:"Payload,omitempty"`
}

func (g *Game) AddCommands(cmds Commands) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, cmd := range cmds {
		g.commands <- cmd
	}
}
func (g *Game) AddCommand(cmd Command) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.commands <- cmd

}

func (g *Game) processCommand(cmd Command) {
	player, ok := g.getPlayer(cmd.PlayerID)
	if !ok {
		return
	}

	if cmd.Action == CommandActionReady {
		g.spawn(player)
		return
	}

	if player.Status != PlayerStatusAlive {
		return
	}

	switch cmd.Action {
	case CommandActionReady:
		g.spawn(player)
	case CommandActionMove:
		g.handlePlayerMove(player, cmd.Payload.(string))
	case CommandActionRotate:
		var dir Direction
		if dir, ok = cmd.Payload.(Direction); !ok {
			dir = Direction(cmd.Payload.(int64))
		}
		g.handlePlayerRotate(player, dir)
	case CommandActionBoost:
		g.handlePlayerBoost(player, cmd.Payload.(bool))
	case CommandActionBrake:
		g.handlePlayerBrake(player)
	case CommandActionTeleport:
		g.handlePlayerTeleport(player)
	case CommandActionBlink:
		g.handlePlayerBlink(player)
	case CommandActionHook:
		g.handlePlayerUseHook(player)

	default:
		slog.Error("unknown command action: %s", cmd.Action)
	}
}

func (g *Game) processCommands() {
	slog.Debug("processing commands")
	processing := true
	go func() {
		defer func() {
			g.mu.Lock()
			processing = false
			g.commandCond.Signal()
			g.mu.Unlock()
		}()
		slog.Debug("total commands", "size", len(g.commands))
		for len(g.commands) > 0 {
			cmd := <-g.commands
			g.processCommand(cmd)
		}

	}()
	for processing {
		g.commandCond.Wait()
	}
	slog.Debug("commands processed", "time", time.Since(g.LastTick).Milliseconds())
}
