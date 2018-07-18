package player

import (
	"net"

	"github.com/wascript3r/winter-is-coming/lib/zombie"
)

// Player struct holds information about client
type Player struct {
	ID      int
	Started bool
	Name    string
	Points  int
	Conn    net.Conn
}

// New returns a pointer to a new player
func New(c net.Conn) *Player {
	return &Player{Conn: c}
}

// Shoot fires at given coordinates
func (p *Player) Shoot(x, y int, z *zombie.Zombie) bool {
	if x == z.X && y == z.Y {
		p.Points++
		return true
	}
	return false
}
