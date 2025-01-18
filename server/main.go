package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tinylib/msgp/msgp"

	"wars/game"
)

func getColors() map[*game.RGBA]bool {
	return map[*game.RGBA]bool{
		game.ColorGray:         false,
		game.ColorBrightRed:    false,
		game.ColorBrightGreen:  false,
		game.ColorBrightBlue:   false,
		game.ColorBrightYellow: false,
		game.ColorAqua:         false,
		game.ColorFuchsia:      false,
		game.ColorMaroon:       false,
		game.ColorGreen:        false,
		game.ColorNavy:         false,
		game.ColorOlive:        false,
		game.ColorTeal:         false,
		game.ColorPurple:       false,
		game.ColorSilver:       false,
		game.ColorOrange:       false,
		game.ColorIndigo:       false,
		game.ColorPink:         false,
		game.ColorBrown:        false,
		game.ColorGold:         false,
		game.ColorYellowGreen:  false,
	}
}

type server struct {
	tcp       *net.TCPListener
	udp       *net.UDPConn
	clients   map[string]*srvClient
	game      *game.Game
	fpsTicker *time.Ticker
	quit      chan struct{}
	colors    map[*game.RGBA]bool

	mu sync.Mutex
}

func main() {
	msgp.RegisterExtension(99, func() msgp.Extension { return new(game.RGBA) })
	msgp.RegisterExtension(100, func() msgp.Extension { return new(game.Direction) })
	msgp.RegisterExtension(101, func() msgp.Extension { return new(game.PlayerStatus) })

	var tcpAddr, udpAddr string
	flag.StringVar(&tcpAddr, "tcpAddr", ":4200", "Server tcp address")
	flag.StringVar(&udpAddr, "udpAddr", ":4201", "Server udp address")
	flag.Parse()

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)

	srv := &server{
		game:      game.NewGame(),
		clients:   make(map[string]*srvClient),
		colors:    getColors(),
		fpsTicker: time.NewTicker(time.Millisecond * 16),
		quit:      make(chan struct{}),
	}

	err := srv.listen(tcpAddr, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.close()

	srv.game.Start()
	srv.initTickers()

	slog.Info("server started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	srv.game.Stop()

	slog.Info("got signal")
	slog.Info("server stopped")
}
