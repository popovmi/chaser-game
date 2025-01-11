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
	
	"wars/lib/colors"
	"wars/lib/game"
	"wars/lib/messages"
)

func getColors() map[colors.RGBA]bool {
	return map[colors.RGBA]bool{
		colors.Green:      false,
		colors.Blue:       false,
		colors.Yellow:     false,
		colors.Purple:     false,
		colors.LightBlue:  false,
		colors.Sky:        false,
		colors.Lime:       false,
		colors.Orange:     false,
		colors.LightGreen: false,
		colors.Brown:      false,
	}
}

type server struct {
	tcp        *net.TCPListener
	udp        *net.UDPConn
	clients    map[string]*srvClient
	game       *game.Game
	rateTicker *time.Ticker
	quit       chan struct{}
	colors     map[colors.RGBA]bool

	mu sync.Mutex
}

func main() {
	msgp.RegisterExtension(98, func() msgp.Extension { return new(messages.MessageBody) })
	msgp.RegisterExtension(99, func() msgp.Extension { return new(colors.RGBA) })

	var tcpAddr, udpAddr string
	flag.StringVar(&tcpAddr, "tcpAddr", ":4200", "Server tcp address")
	flag.StringVar(&udpAddr, "udpAddr", ":4201", "Server udp address")
	flag.Parse()

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)

	srv := &server{
		game:       game.NewGame(),
		clients:    make(map[string]*srvClient),
		colors:     getColors(),
		rateTicker: time.NewTicker(time.Second / 60),
		quit:       make(chan struct{}),
	}

	err := srv.listen(tcpAddr, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.close()

	srv.initTickers()
	slog.Info("tickers started")
	slog.Info("server started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	slog.Info("got signal")
	slog.Info("server stopped")
}
