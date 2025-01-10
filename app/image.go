package main

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"wars/lib/game"
)

var (
	//go:embed assets/background.png
	bgBytes []byte

	//go:embed assets/brick.png
	brickBytes []byte

	//go:embed assets/portal.png
	portalBytes []byte
)

const faceLength = 30

func (c *gameClient) createDefaultImages() {
	background, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(bgBytes))
	if err != nil {
		panic(err)
	}
	brick, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(brickBytes))
	if err != nil {
		panic(err)
	}
	portal, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(portalBytes))
	if err != nil {
		panic(err)
	}

	worldImg := ebiten.NewImage(game.FieldWidth, game.FieldHeight)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.35)
	worldImg.DrawImage(background, op)

	portalRealLength := float64(2 * game.PortalRadius)
	portalImageLength := float64(portal.Bounds().Dx())
	portalScale := portalRealLength / portalImageLength
	portalOp := &ebiten.DrawImageOptions{}
	portalOp.GeoM.Scale(portalScale, portalScale)
	portalImg := ebiten.NewImage(int(portalRealLength), int(portalRealLength))
	portalImg.DrawImage(portal, portalOp)

	brickImg := ebiten.NewImage(200, 40)
	brickOp := &ebiten.DrawImageOptions{}
	brickOp.ColorScale.ScaleAlpha(0.70)
	brickImg.DrawImage(brick, brickOp)

	c.worldImg = worldImg
	c.portalImg = portalImg
	c.brickImg = brickImg
}

func (c *gameClient) CreatePlayerImages(p *game.Player) {
	w, h := 2*game.Radius, 2*game.Radius

	baseImg := ebiten.NewImage(w, h)
	vector.DrawFilledCircle(baseImg, float32(game.Radius), float32(game.Radius), game.Radius, p.Color.ToColorRGBA(),
		true)
	vector.StrokeLine(baseImg,
		float32(game.Radius)+faceLength, float32(game.Radius),
		float32(game.Radius), float32(game.Radius),
		5, color.Black, true,
	)
	drawEyes(baseImg, float32(game.Radius), float32(game.Radius))

	chaseImg := ebiten.NewImage(w, h)
	vector.StrokeCircle(chaseImg, float32(game.Radius), float32(game.Radius), game.Radius-2, 3,
		p.Color.ToColorRGBA(),
		true)
	vector.StrokeLine(chaseImg,
		float32(game.Radius)+faceLength, float32(game.Radius),
		float32(game.Radius), float32(game.Radius),
		5, p.Color.ToColorRGBA(), true,
	)
	drawEyes(chaseImg, float32(game.Radius), float32(game.Radius))

	c.playerImages[p.ID] = &playerImg{baseImg, chaseImg}
}

func drawEyes(img *ebiten.Image, cx, cy float32) {
	eyeRadius := float32(game.Radius) / 3.5
	apRadius := float32(game.Radius) / 6.0
	vector.DrawFilledCircle(img, cx, cy-eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy-eyeRadius*1.3, apRadius, color.Black, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius*1.3, apRadius, color.Black, true)
}
