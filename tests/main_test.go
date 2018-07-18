package tests

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/wascript3r/winter-is-coming/server"
)

func init() {
	go func() {
		server.Run()
	}()
}

func newConn(t *testing.T) (net.Conn, *bufio.Scanner) {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Fatal(err)
	}
	return conn, bufio.NewScanner(conn)
}

func skipScan(t *testing.T, sc *bufio.Scanner, n int) {
	for i := 0; i < n; i++ {
		require.True(t, sc.Scan())
	}
}

func skipHelpMsg(t *testing.T, sc *bufio.Scanner) {
	skipScan(t, sc, 7)
}

func wait() {
	time.Sleep(500 * time.Millisecond)
}
