package main

import (
	"log"
	"net"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinylib/msgp/msgp"

	"wars/app/components"
	"wars/lib/colors"
	"wars/lib/game"
)

const (
	defaultWindowWidth  = 1200
	defaultWindowHeight = 900
)

type gameScreen = byte

const (
	screenWait gameScreen = iota
	screenMain
	screenGame
)

type playerImg struct {
	animation *Animation
	baseImg   *ebiten.Image
	hookImg   *ebiten.Image
	astroImg  *ebiten.Image
}

type untouchableTimer struct {
	t       int
	visible bool
}

type gameClient struct {
	clientID string

	screen    gameScreen
	nameInput *components.TextField

	game  *game.Game
	audio *music

	windowW          int
	windowH          int
	cameraX, cameraY float64

	worldImg         *ebiten.Image
	portalStaticImg  *ebiten.Image
	portalAnimations map[string]*Animation
	brickImg         *ebiten.Image
	healthImg        *ebiten.Image
	healthFillImg    *ebiten.Image
	playerImg        *ebiten.Image
	playerImages     map[string]*playerImg

	untouchableTimers map[string]*untouchableTimer

	tcpAddr string
	udpAddr string
	TCPConn net.Conn
	UDPConn *net.UDPConn

	mu sync.Mutex
}

var tcpAddr, udpAddr string

func main() {
	msgp.RegisterExtension(99, func() msgp.Extension { return new(colors.RGBA) })

	LoadFonts()

	c := &gameClient{
		audio:             newGameMusic(),
		playerImages:      map[string]*playerImg{},
		untouchableTimers: map[string]*untouchableTimer{},
		screen:            screenWait,
		tcpAddr:           tcpAddr,
		udpAddr:           udpAddr,
	}

	c.createDefaultImages()

	go c.openTCPConnection()
	defer func() {
		c.TCPConn.Close()
		if c.UDPConn != nil {
			c.UDPConn.Close()
		}
	}()

	ebiten.SetWindowTitle("WARS")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
