package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"wars/lib/game"
)

var (
	//go:embed assets/background2.png
	bgBytes []byte

	//go:embed assets/brick.png
	brickBytes []byte

	//go:embed assets/portal.png
	portalBytes []byte

	//go:embed assets/ship/Ship1.png
	ship1Bytes []byte
	//go:embed assets/ship/Ship2.png
	ship2Bytes []byte
	//go:embed assets/ship/Ship3.png
	ship3Bytes []byte
	//go:embed assets/ship/Ship4.png
	ship4Bytes []byte
	//go:embed assets/ship/Ship5.png
	ship5Bytes []byte
	//go:embed assets/ship/Ship6.png
	ship6Bytes []byte
	//go:embed assets/ship/Ship7.png
	ship7Bytes []byte
	//go:embed assets/ship/Ship8.png
	ship8Bytes    []byte
	playerSprites = make([]*ebiten.Image, 0)

	//go:embed assets/portal/Portal_01.png
	portal1Bytes []byte
	//go:embed assets/portal/Portal_02.png
	portal2Bytes []byte
	//go:embed assets/portal/Portal_03.png
	portal3Bytes []byte
	//go:embed assets/portal/Portal_04.png
	portal4Bytes []byte
	//go:embed assets/portal/Portal_05.png
	portal5Bytes []byte
	//go:embed assets/portal/Portal_06.png
	portal6Bytes []byte
	//go:embed assets/portal/Portal_07.png
	portal7Bytes []byte
	//go:embed assets/portal/Portal_08.png
	portal8Bytes []byte
	//go:embed assets/portal/Portal_09.png
	portal9Bytes  []byte
	portalSprites = make([]*ebiten.Image, 0)

	//go:embed assets/Player100x100.png
	astroBytes []byte
)

const faceLength = 30

func (c *gameClient) createDefaultImages() {
	{
		healthImage := ebiten.NewImage(2*game.Radius, 6)
		healthImage.Fill(color.RGBA{R: 128, A: 255})
		c.healthImg = healthImage
	}

	{
		healthFillImg := ebiten.NewImage(2*game.Radius, 6)
		healthFillImg.Fill(color.RGBA{G: 255, A: 255})
		c.healthFillImg = healthFillImg
	}

	{
		background, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(bgBytes))
		if err != nil {
			panic(err)
		}
		worldImg := ebiten.NewImage(game.FieldWidth, game.FieldHeight)
		bgW, bgH := background.Bounds().Dx(), background.Bounds().Dy()
		for i := range 2 {
			for j := range 2 {
				fi := float64(i)
				fj := float64(j)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(float64(game.FieldWidth/2)/float64(bgW), float64(game.FieldHeight/2)/float64(bgH))
				op.GeoM.Translate(fi*float64(game.FieldWidth)/2, fj*float64(game.FieldHeight)/2)
				op.ColorScale.ScaleAlpha(0.35)
				worldImg.DrawImage(background, op)
			}
		}
		c.worldImg = worldImg
	}
	{
		portal, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(portalBytes))
		if err != nil {
			panic(err)
		}
		portalRealLength := float64(2 * game.PortalRadius)
		portalImageLength := float64(portal.Bounds().Dx())
		portalScale := portalRealLength / portalImageLength
		portalOp := &ebiten.DrawImageOptions{}
		portalOp.GeoM.Scale(portalScale, portalScale)
		portalImg := ebiten.NewImage(int(portalRealLength), int(portalRealLength))
		portalImg.DrawImage(portal, portalOp)
		c.portalStaticImg = portalImg
	}
	{
		brick, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(brickBytes))
		if err != nil {
			panic(err)
		}
		brickImg := ebiten.NewImage(200, 40)
		brickOp := &ebiten.DrawImageOptions{}
		brickOp.ColorScale.ScaleAlpha(0.70)
		brickImg.DrawImage(brick, brickOp)
		c.brickImg = brickImg
	}
	{
		astro, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(astroBytes))
		if err != nil {
			panic(err)
		}
		playerImg := ebiten.NewImage(2*game.Radius, 2*game.Radius)
		astroOp := &ebiten.DrawImageOptions{}
		scale := 2 * game.Radius / float64(astro.Bounds().Dx())
		astroOp.GeoM.Scale(scale, scale)
		playerImg.DrawImage(astro, astroOp)
		c.playerImg = playerImg
	}
	for _, spriteBytes := range [][]byte{
		ship1Bytes, ship2Bytes,
		ship3Bytes, ship4Bytes,
		ship5Bytes, ship6Bytes,
		ship7Bytes, ship8Bytes,
	} {
		shipSprite, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(spriteBytes))
		if err != nil {
			panic(err)
		}
		spriteImg := ebiten.NewImage(2*game.Radius, 2*game.Radius)
		op := &ebiten.DrawImageOptions{}
		scale := 2 * game.Radius / float64(shipSprite.Bounds().Dx())
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(-game.Radius, -game.Radius)
		op.GeoM.Rotate(-math.Pi / 2)
		op.GeoM.Translate(game.Radius, game.Radius)
		spriteImg.DrawImage(shipSprite, op)
		playerSprites = append(playerSprites, spriteImg)
	}
	for _, portalBytes := range [][]byte{
		portal1Bytes, portal2Bytes,
		portal3Bytes, portal4Bytes,
		portal5Bytes, portal6Bytes,
		portal7Bytes, portal8Bytes, portal9Bytes,
	} {
		portalSprite, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(portalBytes))
		if err != nil {
			panic(err)
		}
		spriteImg := ebiten.NewImage(2*game.PortalRadius, 2*game.PortalRadius)
		op := &ebiten.DrawImageOptions{}
		scale := 2 * game.PortalRadius / float64(portalSprite.Bounds().Dx())
		op.GeoM.Scale(scale, scale)
		spriteImg.DrawImage(portalSprite, op)
		portalSprites = append(portalSprites, spriteImg)
		c.portalStaticImg = portalSprites[0]
	}
}

func (c *gameClient) createPlayerImages(p *game.Player) {
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

	animation := &Animation{Frames: playerSprites, AnimationSpeed: 0.125, img: playerSprites[0]}

	hookImg := ebiten.NewImage(game.MaxHookLength, 5)
	vector.StrokeLine(
		hookImg,
		0, 0,
		game.MaxHookLength, 5,
		5,
		p.Color.ToColorRGBA(),
		true,
	)

	astroImg := ebiten.NewImage(w, h)
	astroOp := &ebiten.DrawImageOptions{}
	astroOp.ColorScale.ScaleWithColor(p.Color.ToColorRGBA())
	astroImg.DrawImage(c.playerImg, astroOp)

	c.playerImages[p.ID] = &playerImg{animation, baseImg, hookImg, astroImg}
}

func drawEyes(img *ebiten.Image, cx, cy float32) {
	eyeRadius := float32(game.Radius) / 3.5
	apRadius := float32(game.Radius) / 6.0
	vector.DrawFilledCircle(img, cx, cy-eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy-eyeRadius*1.3, apRadius, color.Black, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius*1.3, apRadius, color.Black, true)
}

func (c *gameClient) createPortalsAnimations() {
	c.portalAnimations = map[string]*Animation{}
	for _, portal := range c.game.PortalNetwork.Portals {
		c.portalAnimations[portal.ID] =
			&Animation{Frames: portalSprites, AnimationSpeed: 0.6, img: portalSprites[0], Reversed: true}
	}
}
