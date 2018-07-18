package tests

import (
	"bufio"
	"net"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/wascript3r/winter-is-coming/server"
)

var (
	IP     string
	config = &server.Config{
		BX:         2,
		BY:         2,
		ZombieName: "night-king",
		Port:       3001,
		Interval:   2 * time.Second,
	}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := strconv.Itoa(config.Port)
	IP = "127.0.0.1:" + port
	go func() {
		server.Run(config)
	}()
}

func newConn(t *testing.T) (net.Conn, *bufio.Scanner) {
	conn, err := net.Dial("tcp", IP)
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
