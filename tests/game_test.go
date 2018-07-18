package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wascript3r/winter-is-coming/server"
)

func TestShoot(t *testing.T) {
	t.Parallel()

	conn, sc := newConn(t)

	skipHelpMsg(t, sc)
	fmt.Fprintln(conn, "START test")

	skipScan(t, sc, 1)
	require.Equal(t, "Player test started playing.", sc.Text())

	skipScan(t, sc, 1)
	require.Contains(t, sc.Text(), "WALK")

	f := strings.Fields(sc.Text())
	require.Equal(t, 4, len(f))
	x, y := f[2], f[3]

	fmt.Fprintln(conn, "SHOOT", x, y)

	skipScan(t, sc, 1)
	require.Equal(t, "BOOM test 1 night-king", sc.Text())

	skipScan(t, sc, 1)
	require.Equal(t, "Game ended. Player test won.", sc.Text())

	fmt.Fprintln(conn, "SHOOT", x, y)
	skipScan(t, sc, 1)
	err := "ERR " + server.ErrNotStarted.Error()
	require.Equal(t, err, sc.Text())
}

func TestGameSharing(t *testing.T) {
	t.Parallel()

	conn1, sc1 := newConn(t)
	conn2, sc2 := newConn(t)

	skipHelpMsg(t, sc1)
	fmt.Fprintln(conn1, "SHARE")
	skipScan(t, sc1, 1)
	require.Contains(t, sc1.Text(), "Game ID: ")
	ID := strings.Replace(sc1.Text(), "Game ID: ", "", 1)

	skipHelpMsg(t, sc2)
	fmt.Fprintln(conn2, "JOIN", ID)
	skipScan(t, sc2, 1)
	require.Equal(t, "Connected.", sc2.Text())

	fmt.Fprintln(conn2, "START test2")
	wait()
	fmt.Fprintln(conn1, "START test1")
	wait()

	skipScan(t, sc2, 2)
	require.Equal(t, "Player test1 started playing.", sc2.Text())
}

func TestZombie(t *testing.T) {
	t.Parallel()

	conn, sc := newConn(t)

	skipHelpMsg(t, sc)
	fmt.Fprintln(conn, "START test")

	c := time.After(time.Duration(config.BX*config.BY) * config.Interval)

loop:
	for {
		select {
		case <-c:
			t.Fatal("Zombie timeout")

		default:
			if sc.Scan() && strings.Contains(sc.Text(), "Game ended. Zombie night-king won.") {
				break loop
			}
		}
	}
}
