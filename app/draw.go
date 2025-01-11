package main

import (
	"fmt"
	"image/color"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"wars/lib/colors"
	"wars/lib/game"
)

const lineSpacing = 1.1

func (c *gameClient) Layout(w, h int) (int, int) {
	if c.windowW != w || c.windowH != h {
		c.windowW = w
		c.windowH = h
		if w > game.FieldWidth {
			c.windowW = game.FieldWidth + 100
		}

		if h > game.FieldHeight {
			c.windowH = game.FieldHeight + 100
		}
	}
	return c.windowW, c.windowH
}

func (c *gameClient) Draw(screen *ebiten.Image) {
	c.fps = ebiten.ActualFPS()

	switch c.screen {
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
	combinedText := fmt.Sprintf("%s; %s", fpsText, tpsText)
	w, _ := text.Measure(combinedText, FontFace14, lineSpacing)

	op := &text.DrawOptions{}
	op.LineSpacing = lineSpacing
	op.GeoM.Translate(float64(c.windowW)-w, 0)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, combinedText, FontFace14, op)
}

func (c *gameClient) drawMain(screen *ebiten.Image) {
	label := "Input name:"
	middleW, middleH := float64(c.windowW/2), float64(c.windowH/2)
	textW, textH := text.Measure(label, FontFace22, lineSpacing)
	op := &text.DrawOptions{}
	op.LineSpacing = lineSpacing
	op.GeoM.Translate(middleW-textW, middleH-textH)
	op.ColorScale.ScaleWithColor(colors.White.ToColorRGBA())
	text.Draw(screen, label, FontFace22, op)
	c.nameInput.Draw(screen)
}

func (c *gameClient) drawGame(screen *ebiten.Image) {
	c.drawWorld(screen)
	c.drawSpells(screen)
	c.drawPlayerList(screen)
	c.drawPlayers(screen)
}

func (c *gameClient) drawWorld(screen *ebiten.Image) {
	worldOp := &ebiten.DrawImageOptions{}
	worldOp.GeoM.Translate(-c.cameraX, -c.cameraY)
	screen.DrawImage(c.worldImg, worldOp)
}

func (c *gameClient) drawPlayers(screen *ebiten.Image) {
	for _, p := range c.game.Players {
		image := c.playerImages[p.ID].animation.Image()
		//image := c.playerImages[p.ID].baseImg
		imageW, imageH := float64(image.Bounds().Dx()), float64(image.Bounds().Dy())

		lineSpacing := 1.1
		nameStr := p.Name
		if p.ID == c.clientID {
			nameStr += " (you)"
		}
		textW, textH := text.Measure(nameStr, FontFaceBold18, lineSpacing)

		isDead := p.Status == game.PlayerStatusDead
		timeSinceDeath := time.Since(p.DeadAt).Seconds()
		if isDead {
			if timeSinceDeath > 1 {
				continue
			}
		}

		imageDrawOptions := &ebiten.DrawImageOptions{}
		nameTextOptions := &text.DrawOptions{}

		imageDrawOptions.GeoM.Translate(-imageW/2, -imageH+game.Radius)
		imageDrawOptions.GeoM.Rotate(p.Angle)
		imageDrawOptions.GeoM.Translate(p.Position.X-c.cameraX, p.Position.Y-c.cameraY)

		nameTextOptions.LineSpacing = lineSpacing
		nameTextOptions.GeoM.Translate(-c.cameraX, -c.cameraY)
		nameTextOptions.GeoM.Translate(p.Position.X-textW/2, p.Position.Y+game.Radius)

		if isDead {
			imageDrawOptions.ColorScale.ScaleWithColor(color.Gray{Y: 50})
			imageDrawOptions.ColorScale.ScaleAlpha(float32(1 - timeSinceDeath))

			nameTextOptions.ColorScale.ScaleWithColor(color.Gray{Y: 50})
			nameTextOptions.ColorScale.ScaleAlpha(float32(1 - timeSinceDeath))
		} else {
			if ut, ok := c.untouchableTimers[p.ID]; ok && !ut.visible {
				imageDrawOptions.ColorScale.ScaleAlpha(0)
			}
			if p.Blinking {
				var alpha float32
				progress := time.Since(p.BlinkedAt).Seconds() / game.BlinkDuration
				if progress < 0.5 {
					alpha = float32(1 - 2*progress)
				} else {
					alpha = float32(2*progress - 1)
				}
				imageDrawOptions.ColorScale.ScaleAlpha(alpha)
			}

			nameTextOptions.ColorScale.ScaleWithColor(p.Color.ToColorRGBA())

			hpOp := &ebiten.DrawImageOptions{}
			hpOp.GeoM.Translate(-c.cameraX, -c.cameraY)
			hpOp.GeoM.Translate(p.Position.X-game.Radius, p.Position.Y-game.Radius-textH/2)
			screen.DrawImage(c.healthImg, hpOp)
			hpWidth := (p.HP / game.MaxHP) * 50
			hpOp.GeoM.Reset()
			hpOp.GeoM.Scale(hpWidth/50, 1)
			hpOp.GeoM.Translate(-c.cameraX, -c.cameraY)
			hpOp.GeoM.Translate(p.Position.X-game.Radius, p.Position.Y-game.Radius-textH/2)
			screen.DrawImage(c.healthFillImg, hpOp)
		}
		text.Draw(screen, nameStr, FontFaceBold18, nameTextOptions)
		screen.DrawImage(image, imageDrawOptions)

		if p.Hook != nil {
			shiftedStartX := float32(p.Position.X) - float32(c.cameraX)
			shiftedStartY := float32(p.Position.Y) - float32(c.cameraY)
			shiftedEndX := float32(p.Hook.End.X) - float32(c.cameraX)
			shiftedEndY := float32(p.Hook.End.Y) - float32(c.cameraY)

			vector.StrokeLine(
				screen,
				shiftedStartX, shiftedStartY,
				shiftedEndX, shiftedEndY,
				3,
				p.Color.ToColorRGBA(),
				true,
			)
		}
		if p.ID == c.clientID {
			var touching bool
			for _, link := range c.game.PortalLinks {
				for _, port := range []*game.Portal{link.P1, link.P2} {
					if port.Touching(p) {
						touching = true

						clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
						if time.Since(link.LastUsed[p.ID]).Seconds() < game.PortalCooldown {
							clr = color.RGBA{R: 255, G: 0, B: 0, A: 255}
						}
						vector.StrokeCircle(screen,
							float32(port.Pos.X-c.cameraX),
							float32(port.Pos.Y-c.cameraY),
							game.PortalRadius-5,
							1, clr, true,
						)
						break
					}
				}
				if touching {
					break
				}
			}
		}
	}
}

func (c *gameClient) drawPortals() {
	for _, link := range c.game.PortalLinks {
		for _, p := range []*game.Portal{link.P1, link.P2} {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.Pos.X-game.PortalRadius, p.Pos.Y-game.PortalRadius)
			c.worldImg.DrawImage(c.portalImg, op)
		}
	}
}

func (c *gameClient) drawBricks() {
	for _, brick := range c.game.Bricks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-brick.W/2, -brick.H/2)
		op.GeoM.Rotate(brick.A)
		op.GeoM.Translate(brick.Pos.X+brick.W/2, brick.Pos.Y+brick.H/2)

		img := c.brickImg
		c.worldImg.DrawImage(img, op)
	}
}

func (c *gameClient) drawPlayerList(screen *ebiten.Image) {
	sorted := make([]*game.Player, 0)
	for _, v := range c.game.Players {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].JoinedAt.Before(sorted[j].JoinedAt)
	})

	i := 1
	for _, player := range sorted {
		playerStr := fmt.Sprintf("%s | %d | %d | %dm", player.Name, player.Kills, player.Deaths,
			int(time.Since(player.JoinedAt).Minutes()))
		textW, textH := text.Measure(playerStr, FontFace18, lineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = lineSpacing
		op.GeoM.Translate(float64(c.windowW)-textW, float64(c.windowH)-textH*float64(i))
		op.ColorScale.ScaleWithColor(player.Color.ToColorRGBA())
		text.Draw(screen, playerStr, FontFace18, op)
		i++
	}
}

func (c *gameClient) drawSpells(screen *ebiten.Image) {
	p := c.game.Players[c.clientID]

	blinkText := "Blink:  Space"
	hookText := "Hook:   Q"
	portalText := "Portal: E"
	strafeText := "Brake:  Shift"

	portalTouching, usedTime := c.game.CanUsePortal(p.ID)
	portalUsedTime := time.Since(usedTime).Seconds()
	if portalTouching {
		if portalUsedTime < game.PortalCooldown {
			portalText = fmt.Sprintf("Portal: %ds", int(game.PortalCooldown-portalUsedTime))
		}
	}

	blinksUsedTime := time.Since(p.BlinkedAt).Seconds()
	if blinksUsedTime < game.BlinkCooldown {
		blinkText = fmt.Sprintf("Blink:  %ds", int(game.BlinkCooldown-blinksUsedTime))
	}

	hookUsedTime := time.Since(p.HookedAt).Seconds()
	if hookUsedTime < game.HookCooldown {
		hookText = fmt.Sprintf("Hook:   %ds", int(game.HookCooldown-hookUsedTime))
	}

	i := 0
	for _, str := range []string{blinkText, hookText, portalText, strafeText} {
		_, textH := text.Measure(str, FontFace14, lineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = lineSpacing
		op.GeoM.Translate(0, textH*float64(i))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, str, FontFace14, op)
		i++
	}

}
