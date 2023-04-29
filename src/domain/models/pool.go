package models

type Pool struct {
	Clients map[string]*Client
	Pairs   chan *Client
}

func NewPool() *Pool {
	pool := Pool{
		Clients: make(map[string]*Client),
		Pairs:   make(chan *Client),
	}

	return &pool
}
