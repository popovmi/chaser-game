package main

import (
	"bytes"
	_ "embed"
	"log/slog"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"wars/game"
)

var (
	//go:embed assets/wall.ogg
	SoundWallOGG []byte

	//go:embed assets/portal.ogg
	SoundPortalOGG []byte

	//go:embed assets/touch.ogg
	SoundTouchOGG []byte
)

type gameSound struct {
	source        []byte
	player        *audio.Player
	playsByPlayer map[string]int64

	mu sync.Mutex
}

func newGameSound(context *audio.Context, vol float64, source []byte) *gameSound {
	stream, err := vorbis.DecodeF32(bytes.NewReader(source))
	if err != nil {
		panic(err)
	}

	player, err := context.NewPlayerF32(stream)
	if err != nil {
		panic(err)
	}
	player.SetVolume(vol)

	return &gameSound{source: source, player: player, playsByPlayer: make(map[string]int64)}
}

func (s *gameSound) play(pId string) {
	now := time.Now().UnixMilli()
	if lastPlayed, ok := s.playsByPlayer[pId]; ok {
		if now-lastPlayed < 750 {
			return
		}
	}

	err := s.player.Rewind()
	if err != nil {
		slog.Error("Could not rewind sound", "error", err.Error())
		return
	}
	s.player.Play()
	s.playsByPlayer[pId] = now
}

func (s *gameSound) plays(ids []string) {
	now := time.Now().UnixMilli()
	shouldPlay := false
	for _, id := range ids {
		if lastPlayed, ok := s.playsByPlayer[id]; !ok || now-lastPlayed >= 750 {
			s.playsByPlayer[id] = now
			shouldPlay = true
		}
	}
	if !shouldPlay {
		return
	}
	err := s.player.Rewind()
	if err != nil {
		slog.Error("Could not rewind sound", "error", err.Error())
		return
	}
	s.player.Play()

}

type music struct {
	audioContext *audio.Context
	wallHit      *gameSound
	portal       *gameSound
	touch        *gameSound
	events       chan game.Event
}

func newGameMusic() *music {
	audioContext := audio.NewContext(44100)

	return &music{
		audioContext: audioContext,
		wallHit:      newGameSound(audioContext, 0.2, SoundWallOGG),
		portal:       newGameSound(audioContext, 1, SoundPortalOGG),
		touch:        newGameSound(audioContext, 0.1, SoundTouchOGG),
		events:       make(chan game.Event),
	}
}

func (m *music) playPortal() {
	err := m.portal.player.Rewind()
	if err != nil {
		slog.Error("Could not rewind portal sound", "error", err.Error())
		return
	}
	m.portal.player.Play()
}

func (m *music) ListenerID() string {
	return "music"
}

func (m *music) Chan() chan game.Event {
	return m.events
}

func (m *music) Listen(stop chan struct{}) {
	for {
		select {
		case e, ok := <-m.events:
			if !ok {
				return
			}
			m.handleGameEvent(e)
		case <-stop:
			return
		}
	}
}

func (m *music) handleGameEvent(event game.Event) {

	switch event.Action {
	case game.EventActionPortalUsed:
		m.playPortal()
	case game.EventActionWallCollide:
		p := event.Payload.(*game.Player)
		m.wallHit.play(p.ID)
	case game.EventActionPlayerCollide:
		ids := event.Payload.([]string)
		m.touch.plays(ids)
	default:

	}
}
