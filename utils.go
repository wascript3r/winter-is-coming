package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/wascript3r/winter/lib/game"
)

func emitErr(conn net.Conn, err error) {
	fmt.Fprintln(conn, "ERR "+err.Error())
}

func showHelp(conn net.Conn) {
	fmt.Fprintln(conn, `
		Command list:
		START {player} - starts a new game (ex. START John)
		SHOOT {x} {y} - shoots at given coordinates (ex. SHOOT 0 1)
	`)
}

func convCoord(xs, ys string, g *game.Game) (int, int, error) {
	x, err := strconv.Atoi(xs)
	if err != nil {
		return 0, 0, ErrInvalidCoord
	}

	y, err := strconv.Atoi(ys)
	if err != nil {
		return 0, 0, ErrInvalidCoord
	}

	if x < 0 || y < 0 || x > g.BX-1 || y > g.BY-1 {
		return 0, 0, ErrInvalidCoord
	}

	return x, y, nil
}
