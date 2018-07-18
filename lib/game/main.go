package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/wascript3r/winter-is-coming/lib/player"
	"github.com/wascript3r/winter-is-coming/lib/repeat"
	"github.com/wascript3r/winter-is-coming/lib/rnd"
	"github.com/wascript3r/winter-is-coming/lib/zombie"
)

type Game struct {
	ID     string
	BX, BY int
	Zombie *zombie.Zombie

	players map[int]*player.Player
	inc     int
	end     chan<- struct{}
	mx      *sync.RWMutex
	started bool
}

func New() *Game {
	return &Game{mx: &sync.RWMutex{}}
}

func (g *Game) Init(bX, bY int, z *zombie.Zombie) {
	g.BX = bX
	g.BY = bY
	g.Zombie = z
	g.players = make(map[int]*player.Player)
}

func (g *Game) Share() string {
	g.ID = rnd.String(7)
	return g.ID
}

func (g *Game) Join(p *player.Player) {
	p.Started = true

	g.mx.Lock()
	g.inc++
	ID := g.inc
	g.players[ID] = p
	g.mx.Unlock()

	p.ID = ID
}

func (g *Game) Leave(p *player.Player) {
	delete(g.players, p.ID)
}

func (g *Game) IsEmpty() bool {
	return len(g.players) == 0
}

func (g *Game) IsShared() bool {
	return g.ID != ""
}

func (g *Game) IsStarted() bool {
	return g.started
}

func (g *Game) Start() {
	if g.started {
		return
	}

	g.end = repeat.Do(2*time.Second, func() bool {
		end := g.Zombie.Walk()
		g.EmitAll("WALK", g.Zombie.Name, g.Zombie.X, g.Zombie.Y)
		if end {
			g.EmitAll("Game ended. Zombie", g.Zombie.Name, "won.")
			g.End()
		}
		return end
	}, false)

	g.started = true
}

func (g *Game) End() {
	if !g.started {
		return
	}

	g.mx.RLock()
	for _, p := range g.players {
		p.Started = false
	}
	g.mx.RUnlock()

	close(g.end)
	g.started = false
	g.init()
}

func (g *Game) init() {
	g.players = make(map[int]*player.Player)
	g.Zombie = zombie.New(g.Zombie.Name, g.BX, g.BY)
}

func (g *Game) EmitAll(a ...interface{}) {
	for _, p := range g.players {
		fmt.Fprintln(p.Conn, a...)
	}
}
