package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	warsgame "wars/lib/game"
)

var (
	//go:embed assets/background.png
	bgBytes []byte

	//go:embed assets/brick.png
	brickBytes []byte

	//go:embed assets/portal.png
	portalBytes []byte
)

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

	invisiblePlayerImg := ebiten.NewImage(2*warsgame.Radius, 2*warsgame.Radius)
	vector.DrawFilledCircle(
		invisiblePlayerImg, warsgame.Radius, warsgame.Radius, warsgame.Radius, color.RGBA{}, true,
	)

	worldImg := ebiten.NewImage(warsgame.FieldWidth, warsgame.FieldHeight)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.35)
	worldImg.DrawImage(background, op)

	portalRealLength := float64(2 * warsgame.PortalRadius)
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

	c.ui.worldImg = worldImg
	c.ui.portalImg = portalImg
	c.ui.brickImg = brickImg
	c.ui.invisiblePlayerImg = invisiblePlayerImg
}

func (c *gameClient) createPlayerImages(p *warsgame.Player, clr color.RGBA) {

	nameStr := p.Name
	if p.ID == c.id {
		nameStr += " (you)"
	}
	textW, textH := text.Measure(nameStr, FontFaceBold18, warsgame.LineSpacing)
	w, h := max(2*warsgame.Radius, textW), 2*warsgame.Radius+textH

	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(p.Color.ToColorRGBA())
	textOp.LineSpacing = warsgame.LineSpacing
	textOp.GeoM.Translate(w/2-textW/2, 0)
	iw, ih := int(w), int(h)

	baseImg := ebiten.NewImage(iw, ih)
	vector.DrawFilledCircle(baseImg, float32(w/2), float32(warsgame.Radius+textH), warsgame.Radius, clr, true)
	text.Draw(baseImg, nameStr, FontFaceBold18, textOp)

	chaseImg := ebiten.NewImage(iw, ih)
	vector.StrokeCircle(chaseImg, float32(w/2), float32(warsgame.Radius+textH), warsgame.Radius-5, 5, clr, true)
	text.Draw(chaseImg, nameStr, FontFaceBold18, textOp)

	c.ui.playerImgs[p.ID] = &playerImg{baseImg, chaseImg}
}
