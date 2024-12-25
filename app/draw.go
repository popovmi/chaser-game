package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"sort"
	"time"
	color2 "wars/lib/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	warsgame "wars/lib/game"
)

func (c *gameClient) Layout(w, h int) (int, int) {
	if c.ui.windowW != w || c.ui.windowH != h {
		c.ui.windowW = w
		c.ui.windowH = h

		if w > warsgame.FieldWidth {
			c.ui.windowW = warsgame.FieldWidth + 100
		}

		if h > warsgame.FieldHeight {
			c.ui.windowH = warsgame.FieldHeight + 100
		}

		c.ui.nameInput = nil
	}
	return c.ui.windowW, c.ui.windowH
}

func (c *gameClient) Draw(screen *ebiten.Image) {
	c.fps = ebiten.ActualFPS()

	switch c.ui.screen {
	case screenMain:
		c.drawMain(screen)
	case screenGame:
		c.drawGame(screen)
	default:
	}

	c.drawStats(screen)
}

func (c *gameClient) drawStats(screen *ebiten.Image) {
	fpsText := fmt.Sprintf("FPS: %.2f", c.fps)
	tpsText := fmt.Sprintf("TPS: %.2f", c.tps)
	ping := fmt.Sprintf("Ping: %d", c.ping)
	combinedText := fmt.Sprintf("%s; %s; %s", fpsText, tpsText, ping)
	w, _ := text.Measure(combinedText, FontFace14, warsgame.LineSpacing)

	op := &text.DrawOptions{}
	op.LineSpacing = warsgame.LineSpacing
	op.GeoM.Translate(float64(c.ui.windowW)-w, 0)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, combinedText, FontFace14, op)
}

func (c *gameClient) drawMain(screen *ebiten.Image) {
	label := "Input name:"
	middleW, middleH := float64(c.ui.windowW/2), float64(c.ui.windowH/2)
	textW, textH := text.Measure(label, FontFace22, warsgame.LineSpacing)
	op := &text.DrawOptions{}
	op.LineSpacing = warsgame.LineSpacing
	op.GeoM.Translate(middleW-textW, middleH-textH)
	op.ColorScale.ScaleWithColor(color2.White.ToColorRGBA())
	text.Draw(screen, label, FontFace22, op)
	c.ui.nameInput.Draw(screen)
}

func (c *gameClient) drawGame(screen *ebiten.Image) {
	c.drawWorld(screen)
	c.drawSpells(screen)
	c.drawPlayerList(screen)
	c.drawChaser(screen)
	c.drawPlayers(screen)
}

func (c *gameClient) drawWorld(screen *ebiten.Image) {
	worldOp := &ebiten.DrawImageOptions{}
	worldOp.GeoM.Translate(-c.ui.cameraX, -c.ui.cameraY)
	screen.DrawImage(c.ui.worldImg, worldOp)
}

func (c *gameClient) drawPlayers(screen *ebiten.Image) {
	for _, p := range c.game.Players {
		img := c.ui.playerImgs[p.ID].baseImg
		if p.ID == c.game.CId {
			img = c.ui.playerImgs[p.ID].chaseImg
		} else if ut, ok := c.ui.untouchableTimers[p.ID]; ok {
			if !ut.visible {
				img = c.ui.invisiblePlayerImg
			}
		}

		w, h := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-c.ui.cameraX, -c.ui.cameraY)
		op.GeoM.Translate(p.Pos.X-w/2, p.Pos.Y-h+warsgame.Radius)
		screen.DrawImage(img, op)

		if p.Hook != nil {
			c.drawHook(screen, p)
		}

		if p.ID == c.id {
			for _, link := range c.game.PortalLinks {
				for _, port := range []*warsgame.Portal{link.P1, link.P2} {
					if port.Touching(p) {
						clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
						if time.Now().UnixMilli()-link.LastUsed[p.ID] < warsgame.PortalCooldown {
							clr = color.RGBA{R: 255, G: 0, B: 0, A: 255}
						}
						vector.StrokeCircle(screen,
							float32(port.Pos.X-c.ui.cameraX),
							float32(port.Pos.Y-c.ui.cameraY),
							warsgame.PortalRadius-5,
							1, clr, true,
						)
					}
				}
			}
		}
	}
}

func (c *gameClient) drawHook(screen *ebiten.Image, p *warsgame.Player) {
	shiftedStartX := float32(p.Pos.X) - float32(c.ui.cameraX)
	shiftedStartY := float32(p.Pos.Y) - float32(c.ui.cameraY)
	shiftedEndX := float32(p.Hook.End.X) - float32(c.ui.cameraX)
	shiftedEndY := float32(p.Hook.End.Y) - float32(c.ui.cameraY)

	vector.StrokeLine(
		screen,
		shiftedStartX, shiftedStartY,
		shiftedEndX, shiftedEndY,
		5,
		p.Color.ToColorRGBA(),
		true,
	)
}

func (c *gameClient) drawPortals() {
	for _, link := range c.game.PortalLinks {
		for _, p := range []*warsgame.Portal{link.P1, link.P2} {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.Pos.X-warsgame.PortalRadius, p.Pos.Y-warsgame.PortalRadius)
			c.ui.worldImg.DrawImage(c.ui.portalImg, op)
		}
	}
}

func (c *gameClient) drawBricks() {
	for _, brick := range c.game.Bricks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-brick.W/2, -brick.H/2)                       // Сдвигаем к центру изображения
		op.GeoM.Rotate(brick.A)                                         // Поворачиваем
		op.GeoM.Translate(brick.Pos.X+brick.W/2, brick.Pos.Y+brick.H/2) // Сдвигаем на позицию кирпича

		img := c.ui.brickImg
		c.ui.worldImg.DrawImage(img, op)
	}
}

func (c *gameClient) drawChaser(screen *ebiten.Image) {
	chaser := c.game.Players[c.game.CId]
	label := fmt.Sprintf("Chaser: %s", chaser.Name)
	textW, textH := text.Measure(label, FontFace18, warsgame.LineSpacing)
	middleW, middleH := float64(c.ui.windowW/2), 25.0
	textOp := &text.DrawOptions{}
	textOp.LineSpacing = warsgame.LineSpacing
	textOp.GeoM.Translate(middleW-textW/2, middleH-textH)
	textOp.ColorScale.ScaleWithColor(chaser.Color.ToColorRGBA())
	text.Draw(screen, label, FontFace18, textOp)
}

func (c *gameClient) drawPlayerList(screen *ebiten.Image) {

	sorted := make([]*warsgame.Player, 0)
	for _, v := range c.game.Players {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].JoinedAt > sorted[j].JoinedAt
	})

	t := time.Now().UnixMilli()
	i := 1
	for _, player := range sorted {
		playerStr := fmt.Sprintf("%s | %dm | %d", player.Name, (t-player.JoinedAt)/1000/60, player.ChaseCount)
		textW, textH := text.Measure(playerStr, FontFace18, warsgame.LineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = warsgame.LineSpacing
		op.GeoM.Translate(float64(c.ui.windowW)-textW, float64(c.ui.windowH)-textH*float64(i))
		op.ColorScale.ScaleWithColor(player.Color.ToColorRGBA())
		text.Draw(screen, playerStr, FontFace18, op)
		i++
	}
}

func (c *gameClient) drawSpells(screen *ebiten.Image) {
	p := c.game.Players[c.id]

	brakeText := "Brake:  Shift"
	blinkText := "Blink:  Space"
	hookText := "Hook:   Q"
	portalText := "Portal: E"

	now := time.Now().UnixMilli()

	portalTouching, usedTime := c.game.CanUsePortal(c.id)
	portalCdLeft := now - usedTime
	if portalTouching {
		if now-usedTime < warsgame.PortalCooldown {
			portalText = fmt.Sprintf("Portal: %ds", (warsgame.PortalCooldown-portalCdLeft)/1000)
		}
	}

	blinkCDLeft := now - p.BlinkedAt
	if blinkCDLeft < warsgame.BlinkCooldown {
		blinkText = fmt.Sprintf("Blink:  %ds", (warsgame.BlinkCooldown-blinkCDLeft)/1000)
	}

	hookCDLeft := now - p.HookedAt
	if hookCDLeft < warsgame.HookCooldown {
		hookText = fmt.Sprintf("Hook:   %ds", (warsgame.HookCooldown-hookCDLeft)/1000)
	}

	i := 0
	for _, str := range []string{brakeText, blinkText, hookText, portalText} {
		_, textH := text.Measure(str, FontFace14, warsgame.LineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = warsgame.LineSpacing
		op.GeoM.Translate(0, textH*float64(i))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, str, FontFace14, op)
		i++
	}

}
