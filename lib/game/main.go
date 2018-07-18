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

// Game struct holds all necessary information about board, zombie, players, etc.
type Game struct {
	ID       string
	BX, BY   int
	Zombie   *zombie.Zombie
	Interval time.Duration

	players map[int]*player.Player
	inc     int
	end     chan<- struct{}
	mx      *sync.RWMutex
	started bool
}

// New returns a pointer to a new game
func New(i time.Duration) *Game {
	return &Game{
		Interval: i,
		mx:       &sync.RWMutex{},
	}
}

// Init initializes created game without additional memory allocation
func (g *Game) Init(bX, bY int, z *zombie.Zombie) {
	g.BX = bX
	g.BY = bY
	g.Zombie = z
	g.players = make(map[int]*player.Player)
}

// Share generates ID for the game in order to be accessed by other clients
func (g *Game) Share() string {
	g.ID = rnd.String(7)
	return g.ID
}

// Join adds new player to connected clients list
func (g *Game) Join(p *player.Player) {
	p.Started = true

	g.mx.Lock()
	g.inc++
	ID := g.inc
	g.players[ID] = p
	g.mx.Unlock()

	p.ID = ID
}

// Leave deletes player from connected clients list
func (g *Game) Leave(p *player.Player) {
	delete(g.players, p.ID)
}

// IsEmpty checks if game still has active clients
func (g *Game) IsEmpty() bool {
	return len(g.players) == 0
}

// IsShared checks if game is accessible by other clients
func (g *Game) IsShared() bool {
	return g.ID != ""
}

// IsStarted checks if game is started
func (g *Game) IsStarted() bool {
	return g.started
}

// Start spawns a new zombie to the game
func (g *Game) Start() {
	if g.started {
		return
	}

	g.end = repeat.Do(g.Interval, func() bool {
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

// End kicks all clients from current game
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

// EmitAll sends a message to all connected clients of the specific game
func (g *Game) EmitAll(a ...interface{}) {
	for _, p := range g.players {
		fmt.Fprintln(p.Conn, a...)
	}
}
