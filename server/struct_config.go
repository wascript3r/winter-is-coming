package server

import (
	"time"
)

type Config struct {
	BX, BY     int
	ZombieName string
	Port       int
	Interval   time.Duration
}
