package main

import (
	"log"
	"net"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinylib/msgp/msgp"

	"chaser/app/components"
	"chaser/lib/colors"
	"chaser/lib/game"
	"chaser/lib/messages"
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
	baseImg  *ebiten.Image
	chaseImg *ebiten.Image
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

	fps float64
	tps float64

	worldImg          *ebiten.Image
	portalImg         *ebiten.Image
	brickImg          *ebiten.Image
	playerImages      map[string]*playerImg
	untouchableTimers map[string]*untouchableTimer

	tcpAddr string
	udpAddr string
	TCPConn net.Conn
	UDPConn *net.UDPConn

	mu sync.Mutex
}

var tcpAddr, udpAddr string

func main() {
	msgp.RegisterExtension(98, func() msgp.Extension { return new(messages.MessageBody) })
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

	ebiten.SetWindowTitle("Chaser")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
