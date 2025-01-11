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
	//go:embed assets/background.png
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
	ship8Bytes []byte

	playerSprites = make([]*ebiten.Image, 0)
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

	animation := &Animation{Frames: playerSprites, AnimationSpeed: 0.125, img: playerSprites[0]}
	c.playerImages[p.ID] = &playerImg{animation, baseImg}
}

func drawEyes(img *ebiten.Image, cx, cy float32) {
	eyeRadius := float32(game.Radius) / 3.5
	apRadius := float32(game.Radius) / 6.0
	vector.DrawFilledCircle(img, cx, cy-eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius, eyeRadius, color.White, true)
	vector.DrawFilledCircle(img, cx, cy-eyeRadius*1.3, apRadius, color.Black, true)
	vector.DrawFilledCircle(img, cx, cy+eyeRadius*1.3, apRadius, color.Black, true)
}
