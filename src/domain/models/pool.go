package models

import "time"

type Pool struct {
	Broadcast chan Message
	Match     chan Message
	Unmatch   chan Message
	Clients   map[string]*Client
	Pairs     []*Client
	CreatedAt time.Time
}

func NewPool() *Pool {
	pool := Pool{
		Broadcast: make(chan Message),
		Match:     make(chan Message),
		Unmatch:   make(chan Message),
		Clients:   make(map[string]*Client),
		Pairs:     make([]*Client, 0),
		CreatedAt: time.Now(),
	}

	return &pool
}
