package server

import (
	"time"
)

// Config holds necessary information for server and game board
type Config struct {
	BX, BY     int
	ZombieName string
	Port       int
	Interval   time.Duration
}
