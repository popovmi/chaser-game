package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Frames            []*ebiten.Image
	CurrentFrameIndex int
	Count             float64
	AnimationSpeed    float64
	Reversed          bool
	direction         int8

	img *ebiten.Image
}

func (a *Animation) Update() {
	if a.direction >= 0 {
		a.Count += a.AnimationSpeed
	} else {
		a.Count -= a.AnimationSpeed
	}
	a.CurrentFrameIndex = int(math.Floor(a.Count))
	if !a.Reversed {
		if a.CurrentFrameIndex >= len(a.Frames) { // restart animation
			a.Count = 0
			a.CurrentFrameIndex = 0
		}
	} else {
		if a.CurrentFrameIndex >= len(a.Frames) { // restart animation
			a.direction = -1
			a.CurrentFrameIndex = len(a.Frames) - 1
		}
		if a.CurrentFrameIndex < 0 {
			a.CurrentFrameIndex = 0
			a.direction = 1
		}
	}

	a.img = a.Frames[a.CurrentFrameIndex]
}

func (a *Animation) Image() *ebiten.Image {
	return a.img
}

func (a *Animation) Reset() {
	a.Count = 0
	a.CurrentFrameIndex = 0
	a.img = a.Frames[0]
}
