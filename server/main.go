package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/wascript3r/winter-is-coming/lib/game"
	"github.com/wascript3r/winter-is-coming/lib/zombie"

	"github.com/wascript3r/winter-is-coming/lib/player"
)

var (
	ErrConfigNotProvided = errors.New("config not provided")
	ErrEmptyCMD          = errors.New("empty command")
	ErrInvalidCMD        = errors.New("invalid command")
	ErrMissingParams     = errors.New("missing command parameters")
	ErrAlreadyStarted    = errors.New("game was already started")
	ErrNotStarted        = errors.New("you must start the game first")
	ErrInvalidCoord      = errors.New("invalid coordinates")
	ErrPleaseWait        = errors.New("please wait until zombie starts walking")
	ErrCannotJoin        = errors.New("cannot join because game is already started")
	ErrGameNotFound      = errors.New("game not found")

	config *Config
	shared = make(map[string]*game.Game)
	mx     = &sync.RWMutex{}
)

// Run starts the server
func Run(c *Config) {
	config = c

	port := strconv.Itoa(config.Port)
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	log.Println("Starting server on port " + port + "...")

	// repeat.Do(time.Second, func() bool {
	// 	log.Println(runtime.NumGoroutine())
	// 	return false
	// }, true)

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	showHelp(conn)

	g := game.New(config.Interval)
	p := player.New(conn)

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		router(&g, p, f)
	}

	g.Leave(p)
	g.EmitAll("Player", p.Name, "disconnected.")
	cleanup(g)
	fmt.Println(shared)
}

// memory leak prevention
func cleanup(g *game.Game) {
	if g.IsEmpty() {
		g.End()
		if g.IsShared() {
			mx.Lock()
			delete(shared, g.ID)
			mx.Unlock()
		}
	}
}

func createGame(g *game.Game, p *player.Player) {
	if g.IsStarted() {
		g.Join(p)
		return
	}
	z := zombie.New(config.ZombieName, config.BX, config.BY)
	g.Init(config.BX, config.BY, z)
	g.Join(p)
	g.Start()
}

func router(gp **game.Game, p *player.Player, f []string) {
	if len(f) == 0 {
		emitErr(p.Conn, ErrEmptyCMD)
		return
	}

	g := *gp

	switch strings.ToLower(f[0]) {
	case "start":
		if len(f) < 2 {
			emitErr(p.Conn, ErrMissingParams)
			return
		}

		if p.Started {
			emitErr(p.Conn, ErrAlreadyStarted)
			return
		}

		name := strings.Join(f[1:], " ")
		p.Name = name
		createGame(g, p)

		g.EmitAll("Player", name, "started playing.")

	case "shoot":
		if len(f) < 3 {
			emitErr(p.Conn, ErrMissingParams)
			return
		}

		if !p.Started {
			emitErr(p.Conn, ErrNotStarted)
			return
		}

		if g.Zombie.NotStarted() {
			emitErr(p.Conn, ErrPleaseWait)
			return
		}

		x, y, err := convCoord(f[1], f[2], g)
		if err != nil {
			emitErr(p.Conn, err)
			return
		}

		if p.Shoot(x, y, g.Zombie) {
			g.EmitAll("BOOM", p.Name, p.Points, g.Zombie.Name)
			g.EmitAll("Game ended. Player", p.Name, "won.")
			g.End()
		} else {
			g.EmitAll("BOOM", p.Name, p.Points)
		}

	case "share":
		if !g.IsShared() {
			g.Share()
			mx.Lock()
			shared[g.ID] = g
			mx.Unlock()
		}
		fmt.Fprintln(p.Conn, "Game ID:", g.ID)

	case "join":
		if len(f) < 2 {
			emitErr(p.Conn, ErrMissingParams)
			return
		}

		if p.Started {
			emitErr(p.Conn, ErrCannotJoin)
			return
		}

		v, ok := shared[f[1]]
		if !ok {
			emitErr(p.Conn, ErrGameNotFound)
			return
		}
		*gp = v
		fmt.Fprintln(p.Conn, "Connected.")

	default:
		emitErr(p.Conn, ErrInvalidCMD)
		showHelp(p.Conn)
	}
}
