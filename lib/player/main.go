package player

import (
	"net"

	"github.com/wascript3r/winter-is-coming/lib/zombie"
)

type Player struct {
	ID      int
	Started bool
	Name    string
	Points  int
	Conn    net.Conn
}

func New(c net.Conn) *Player {
	return &Player{Conn: c}
}

func (p *Player) Shoot(x, y int, z *zombie.Zombie) bool {
	if x == z.X && y == z.Y {
		p.Points++
		return true
	}
	return false
}
