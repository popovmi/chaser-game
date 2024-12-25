package main

import (
	"log"
	"net"
	"sync"

	"wars/app/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinylib/msgp/msgp"

	warscolor "wars/lib/color"
	"wars/lib/game"
	"wars/lib/messages"
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

type gameUI struct {
	windowW int
	windowH int

	screen gameScreen

	nameInput *components.TextField

	cameraX, cameraY float64

	worldImg           *ebiten.Image
	portalImg          *ebiten.Image
	brickImg           *ebiten.Image
	invisiblePlayerImg *ebiten.Image
	playerImgs         map[string]*playerImg

	untouchableTimers map[string]*untouchableTimer
}

type gameClient struct {
	id string

	game  *warsgame.Game
	ui    *gameUI
	audio *music

	fps float64
	tps float64

	tcpAddr string
	udpAddr string
	TCPConn net.Conn
	UDPConn *net.UDPConn
	ping    int

	quit chan struct{}
	mu   sync.Mutex
}

var tcpAddr, udpAddr string

func main() {
	msgp.RegisterExtension(98, func() msgp.Extension { return new(messages.MessageBody) })
	msgp.RegisterExtension(99, func() msgp.Extension { return new(warscolor.RGBA) })

	InitFont()
	c := &gameClient{
		ui: &gameUI{
			screen:            screenWait,
			untouchableTimers: make(map[string]*untouchableTimer),
			playerImgs:        make(map[string]*playerImg),
		},
		audio:   newGameMusic(),
		tcpAddr: tcpAddr,
		udpAddr: udpAddr,
		quit:    make(chan struct{}),
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
	ebiten.SetTPS(warsgame.TPS)
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
