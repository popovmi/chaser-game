package game

import (
	"log/slog"
	"sync"
	"time"
)

//go:generate msgp

type EventListener interface {
	ID() string
	Chan() chan Event
	Listen(chan struct{})
}

type Game struct {
	State *State `msg:"s" json:"State"`

	commands    chan Command
	commandCond sync.Cond

	events    chan Event
	listeners map[string]EventListener

	rateTicker *time.Ticker
	LastTick   time.Time

	stop chan struct{}
	mu   sync.Mutex
	lmu  sync.Mutex
	wg   sync.WaitGroup
}

func NewGame() *Game {

	g := &Game{
		State: NewState(),

		rateTicker: time.NewTicker(time.Millisecond * 16),
		commands:   make(chan Command, 100),
		events:     make(chan Event, 100),
		listeners:  make(map[string]EventListener),

		stop: make(chan struct{}),
	}
	g.commandCond = sync.Cond{L: &g.mu}
	return g
}

func (g *Game) getPlayer(id string) (*Player, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.State.Players[id]
	return p, ok
}

func (g *Game) DeletePlayer(id string) {
	g.mu.Lock()
	delete(g.State.Players, id)
	g.mu.Unlock()
}

func (g *Game) Start() {
	g.wg.Add(2)
	go g.startTicking()
	go g.publishEvents()
	slog.Info("game started")
}

func (g *Game) Stop() {
	close(g.stop)
	close(g.events)
	close(g.commands)
	g.wg.Wait()
	slog.Info("game stopped")
}

func (g *Game) startTicking() {
	defer func() {
		slog.Info("stop ticker")
		g.rateTicker.Stop()
		g.wg.Done()
	}()
	g.LastTick = time.Now()

	for {
		select {
		case <-g.rateTicker.C:
			g.tick()
		case <-g.stop:
			return
		default:
		}
	}
}

func (g *Game) tick() {
	g.mu.Lock()
	defer g.mu.Unlock()

	t := time.Since(g.LastTick).Milliseconds()
	g.LastTick = time.Now()

	slog.Debug("starting tick", "timeSinceLastTick", t)

	dt := float64(t) / 1000

	g.processCommands()
	g.update(dt)

	slog.Debug("ticked", "time", time.Since(g.LastTick).Milliseconds())
}

func (g *Game) Lock() {
	g.mu.Lock()
}

func (g *Game) Unlock() {
	g.mu.Unlock()
}
