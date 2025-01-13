package main

import (
	"fmt"
	"image/color"
	"maps"
	"math"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	switch c.screen {
	case screenMain:
		c.drawMain(screen)
	case screenGame:
		c.drawGame(screen)
		if p, ok := c.game.Players[c.clientID]; ok {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %0.2f, Y: %0.2f", p.Position.X, p.Position.Y), 0,
				c.windowH-20)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %0.2f", p.Velocity.Length()), 0, c.windowH-40)
		}
	default:
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), 0, 20)

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
	c.drawPortals(screen)
	c.drawSpells(screen)
	c.drawPlayers(screen)
	c.drawPlayerList(screen)
}

func (c *gameClient) drawWorld(screen *ebiten.Image) {
	worldOp := &ebiten.DrawImageOptions{}
	worldOp.GeoM.Translate(-c.cameraX, -c.cameraY)
	screen.DrawImage(c.worldImg, worldOp)
}

func (c *gameClient) drawPlayers(screen *ebiten.Image) {
	for _, p := range c.game.Players {
		playerImages, ok := c.playerImages[p.ID]
		if !ok {
			continue
		}
		image := playerImages.astroImg
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
			imageDrawOptions.ColorScale.ScaleWithColor(color.Gray{Y: 150})
			imageDrawOptions.ColorScale.ScaleAlpha(float32(1 - timeSinceDeath))

			nameTextOptions.ColorScale.ScaleWithColor(color.Gray{Y: 150})
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
			if p.Teleporting {
				var alpha float32
				progress := time.Since(p.TeleportedAt).Seconds() / game.TeleportDuration
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
			hpWidth := p.HP / game.MaxHP
			hpOp.GeoM.Reset()
			hpOp.GeoM.Scale(hpWidth, 1)
			hpOp.GeoM.Translate(-c.cameraX, -c.cameraY)
			hpOp.GeoM.Translate(p.Position.X-game.Radius, p.Position.Y-game.Radius-textH/2)
			screen.DrawImage(c.healthFillImg, hpOp)
		}
		text.Draw(screen, nameStr, FontFaceBold18, nameTextOptions)
		screen.DrawImage(image, imageDrawOptions)

		if p.Hook != nil {
			hookImg := c.playerImages[p.ID].hookImg
			hookOp := &ebiten.DrawImageOptions{}
			hookLength := p.HookLength() / game.MaxHookLength
			hookOp.GeoM.Scale(hookLength, 1)
			angle := p.Angle
			if p.Hook.Stuck {
				angle = math.Atan2(p.Hook.End.Y-p.Position.Y, p.Hook.End.X-p.Position.X)
			}
			hookOp.GeoM.Rotate(angle)
			hookOp.GeoM.Translate(p.Position.X-c.cameraX, p.Position.Y-c.cameraY)
			screen.DrawImage(hookImg, hookOp)
		}
	}
}

func (c *gameClient) drawPortals(screen *ebiten.Image) {
	for _, p := range c.game.PortalNetwork.Portals {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.Pos.X-game.PortalRadius-c.cameraX, p.Pos.Y-game.PortalRadius-c.cameraY)

		img := c.portalStaticImg
		t := time.Since(p.LastUsedAt).Seconds()
		if t <= 1.0 {
			img = c.portalAnimations[p.ID].Image()
		}
		screen.DrawImage(img, op)
		player := c.game.Players[c.clientID]
		if !player.Teleporting {
			canUse, portal, _ := c.game.PortalNetwork.CanUsePortal(player)
			if portal != nil && portal.ID == p.ID {
				clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
				if !canUse {
					clr = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				}
				vector.StrokeCircle(screen,
					float32(portal.Pos.X-c.cameraX),
					float32(portal.Pos.Y-c.cameraY),
					game.PortalRadius-5,
					1, clr, true,
				)
			}
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
	for v := range maps.Values(c.game.Players) {
		sorted = append(sorted, v)
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Kills != sorted[j].Kills {
			return sorted[i].Kills > sorted[j].Kills
		}
		if sorted[i].Deaths != sorted[j].Deaths {
			return sorted[i].Deaths < sorted[j].Deaths
		}
		return sorted[i].JoinedAt.Before(sorted[j].JoinedAt)

	})

	i := 0
	for _, player := range sorted {
		playerStr := fmt.Sprintf("%s | %d | %d | %dm", player.Name, player.Kills, player.Deaths,
			int(time.Since(player.JoinedAt).Minutes()))
		textW, textH := text.Measure(playerStr, FontFace18, lineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = lineSpacing
		op.GeoM.Translate(float64(c.windowW)-textW, textH*float64(i))
		op.ColorScale.ScaleWithColor(player.Color.ToColorRGBA())
		text.Draw(screen, playerStr, FontFace18, op)
		i++
	}
}

func (c *gameClient) drawSpells(screen *ebiten.Image) {
	p := c.game.Players[c.clientID]

	blinkText := "Blink: Space"
	hookText := "Hook: Q"
	portalText := "Portal: E"
	brakeText := "Brake: Shift"

	_, _, cooldown := c.game.PortalNetwork.CanUsePortal(p)
	if cooldown != nil && cooldown.Seconds() < game.PortalCooldown {
		portalText = fmt.Sprintf("Portal: %ds", int(game.PortalCooldown-cooldown.Seconds()))
	}

	blinksUsedTime := time.Since(p.BlinkedAt).Seconds()
	if blinksUsedTime < game.BlinkCooldown {
		blinkText = fmt.Sprintf("Blink:  %ds", int(game.BlinkCooldown-blinksUsedTime))
	}

	hookUsedTime := time.Since(p.HookedAt).Seconds()
	if hookUsedTime < game.HookCooldown {
		hookText = fmt.Sprintf("Hook: %ds", int(game.HookCooldown-hookUsedTime))
	}

	combined := fmt.Sprintf("%s\t%s\t%s\t%s", hookText, blinkText, portalText, brakeText)
	textW, textH := text.Measure(combined, FontFace14, lineSpacing)
	op := &text.DrawOptions{}
	op.LineSpacing = lineSpacing
	op.GeoM.Translate(float64(c.windowW)/2-textW/2, float64(c.windowH)-textH)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, combined, FontFace14, op)
}
