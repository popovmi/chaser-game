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

	img *ebiten.Image
}

func (a *Animation) Update() {
	a.Count += a.AnimationSpeed
	a.CurrentFrameIndex = int(math.Floor(a.Count))

	if a.CurrentFrameIndex >= len(a.Frames) { // restart animation
		a.Count = 0
		a.CurrentFrameIndex = 0
	}

	a.img = a.Frames[a.CurrentFrameIndex]
}

func (a *Animation) Image() *ebiten.Image {
	return a.img
}
