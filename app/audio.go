package main

import (
	"bytes"
	_ "embed"
	"log/slog"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
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

type music struct {
	audioContext *audio.Context
	wallHit      *gameSound
	portal       *gameSound
	touch        *gameSound

	mu sync.Mutex
}

func newGameMusic() *music {
	audioContext := audio.NewContext(44100)

	return &music{
		audioContext: audioContext,
		wallHit:      newGameSound(audioContext, 0.2, SoundWallOGG),
		portal:       newGameSound(audioContext, 1, SoundPortalOGG),
		touch:        newGameSound(audioContext, 0.1, SoundTouchOGG),
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
