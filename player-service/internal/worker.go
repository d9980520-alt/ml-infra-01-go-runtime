package internal

import "time"

type Player struct {
	Id       int
	Name     string
	Balance  int64
	LastSeen time.Time
}
