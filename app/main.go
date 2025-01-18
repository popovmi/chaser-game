package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinylib/msgp/msgp"

	"wars/app/components"
	"wars/game"
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

	tcpAddr      string
	TCPConn      net.Conn
	pingInterval time.Duration
	ping         time.Duration
	lastPingTime time.Time

	udpAddr       string
	UDPConn       *net.UDPConn
	udpMsgCounter atomic.Uint64

	mu sync.Mutex
}

var tcpAddr, udpAddr string

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	msgp.RegisterExtension(99, func() msgp.Extension { return new(game.RGBA) })
	msgp.RegisterExtension(100, func() msgp.Extension { return new(game.Direction) })
	msgp.RegisterExtension(101, func() msgp.Extension { return new(game.PlayerStatus) })

	LoadFonts()

	c := &gameClient{
		game:              game.NewGame(),
		audio:             newGameMusic(),
		playerImages:      map[string]*playerImg{},
		untouchableTimers: map[string]*untouchableTimer{},
		screen:            screenWait,
		tcpAddr:           tcpAddr,
		pingInterval:      time.Second,
		udpAddr:           udpAddr,
	}

	c.createDefaultImages()
	c.game.AppendListener(c.audio)

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
