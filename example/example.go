package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	game2 "wars/game"
)

func main() {

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	g := game2.NewGame()
	g.Start()

	var listeners []listener

	var i float64 = 1
	for ; i <= 10; i++ {
		pid := fmt.Sprintf("p%d", int(i))
		player := &game2.Player{ID: pid}
		l := listener{pid, make(chan game2.Event)}
		listeners = append(listeners, l)
		g.AppendListener(&l)
		go func() {
			time.Sleep(time.Millisecond * 10)
			g.Join(player)
			time.Sleep(time.Millisecond * 10)
			g.AddCommand(game2.Command{PlayerID: pid, Action: game2.CommandActionReady, Payload: nil})
			time.Sleep(time.Millisecond * 50)

			for range 100 {
				sleepAndCommand(pid, g)
			}
		}()
	}

	time.Sleep(2 * time.Second)
	g.Stop()
	for _, l := range listeners {
		close(l.events)
	}
	slog.Info("game over")
}

type listener struct {
	id     string
	events chan game2.Event
}

func (l *listener) Listen(stop chan struct{}) {
	for {
		select {
		case _, ok := <-l.events:
			if !ok {
				return
			}
			//slog.Debug("handled game event", "event", e.Action, "payload", e.Payload, "listener", l.id)
		case <-stop:
			return
		}
	}
}

func (l listener) ID() string {
	return l.id
}

func (l *listener) Chan() chan game2.Event {
	return l.events
}

func sleepAndCommand(playerID string, g *game2.Game) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Список доступных действий
	actions := []game2.CommandAction{
		game2.CommandActionMove,
		game2.CommandActionBlink,
		game2.CommandActionRotate,
		game2.CommandActionBoost,
		game2.CommandActionHook,
		game2.CommandActionBrake,
		game2.CommandActionTeleport,
	}

	// Генерируем случайный индекс действия
	actionIndex := rand.Intn(len(actions))
	action := actions[actionIndex]

	// Генерируем случайный payload, если требуется
	var payload interface{}
	switch action {
	case game2.CommandActionMove:
		directions := []string{"u", "d", "l", "r", "ul", "ur", "dl", "dr"}
		payload = directions[rand.Intn(len(directions))]
	case game2.CommandActionRotate:
		payload = game2.Direction(rand.Intn(3))
	case game2.CommandActionBoost:
		payload = rand.Intn(2) == 0 // Случайный bool
	}

	time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
	g.AddCommand(game2.Command{
		Action:   action,
		PlayerID: playerID,
		Payload:  payload,
	})
}
