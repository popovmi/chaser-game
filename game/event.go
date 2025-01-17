package game

import (
	"fmt"
	"log/slog"
	"sync"
)

//go:generate msgp

type EventAction byte

const (
	EventActionNone EventAction = iota
	EventActionSpawned
	EventActionMoved
	EventActionRotated
	EventActionBraked
	EventActionTeleported
	EventActionBlinked
	EventActionHooked
	EventActionReturnedHook
	EventActionBoosted
)

func (ea EventAction) String() string {
	switch ea {
	case EventActionNone:
		return "None"
	case EventActionSpawned:
		return "Spawned"
	case EventActionMoved:
		return "Moved"
	case EventActionRotated:
		return "Rotated"
	case EventActionBraked:
		return "Braked"
	case EventActionTeleported:
		return "Teleported"
	case EventActionBlinked:
		return "Blinked"
	case EventActionHooked:
		return "Hooked"
	case EventActionReturnedHook:
		return "ReturnedHook"
	case EventActionBoosted:
		return "Boosted"
	default:
		return fmt.Sprintf("Unknown(%d)", ea)
	}
}

func (ea EventAction) MarshalJSON() (data []byte, err error) {
	return []byte(fmt.Sprintf("\"%s\"", ea.String())), nil
}

type Event struct {
	Action  EventAction `msg:"act,omitempty" json:"Action,omitempty"`
	Payload interface{} `msg:"pld,omitempty" json:"Payload,omitempty"`
}

func (g *Game) AppendListener(l EventListener) {
	g.lmu.Lock()
	defer g.lmu.Unlock()
	g.listeners[l.ID()] = l
	go l.Listen(g.stop)
}

func (g *Game) publishEvents() {
	defer func() {
		slog.Info("stop events publishing")
		g.wg.Done()
	}()
	for {
		select {
		case event, ok := <-g.events:
			if !ok {
				return
			}

			g.lmu.Lock()
			slog.Debug("publish event", "event", event)
			wg := sync.WaitGroup{}
			wg.Add(len(g.listeners))
			for _, l := range g.listeners {
				go func() {
					l.Chan() <- event
					wg.Done()
				}()
			}
			wg.Wait()
			g.lmu.Unlock()

		case <-g.stop:
			return
		}
	}
}
