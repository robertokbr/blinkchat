package models

import (
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
)

type Pool struct {
	Broadcast chan Message
	Match     chan Message
	Unmatch   chan Message
	Clients   map[string]*Client
	Pairs     []*Client
}

func NewPool() *Pool {
	pool := Pool{
		Broadcast: make(chan Message),
		Match:     make(chan Message),
		Unmatch:   make(chan Message),
		Clients:   make(map[string]*Client),
		Pairs:     make([]*Client, 0),
	}

	return &pool
}

func (p *Pool) PushMessage(message Message) {
	switch message.Action {
	case enums.BROADCASTING:
		p.Broadcast <- message
		break
	case enums.MATCHING:
		p.Match <- message
		break
	case enums.UNMATCHING:
		p.Unmatch <- message
		break
	default:
		logger.Debugf("Receiving message %v", message)
	}
}
