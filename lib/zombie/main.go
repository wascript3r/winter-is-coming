package zombie

import "math/rand"

// Zombie struct holds coordinates and name of the zombie
type Zombie struct {
	BX, BY, X, Y int
	Name         string
}

// New returns a pointer to a new zombie
func New(n string, bX, bY int) *Zombie {
	return &Zombie{
		BX:   bX,
		BY:   bY,
		Name: n,
	}
}

// Walk moves zombie towards the archer
func (z *Zombie) Walk() bool {
	n := rand.Intn(2)

	if n == 0 {
		if z.X < z.BX-1 {
			z.X++
		} else if z.Y < z.BY-1 {
			z.Y++
		}
	} else if n == 1 {
		if z.Y < z.BY-1 {
			z.Y++
		} else if z.X < z.BX-1 {
			z.X++
		}
	}

	return z.X == z.BX-1 && z.Y == z.BY-1
}

// NotStarted returns if zombie is still at initial coordinates
func (z *Zombie) NotStarted() bool {
	return z.X == 0 && z.Y == 0
}
