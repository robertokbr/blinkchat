package models

import (
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/interfaces"
)

type Client struct {
	*User
	Conn     interfaces.WebsocketConnection
	Pair     *Client
	State    enums.UserState
	PairedAt time.Time
}

func NewClient(user *User, connection interfaces.WebsocketConnection) *Client {
	return &Client{
		User:  user,
		Conn:  connection,
		State: enums.NOT_IN_A_MATCH,
	}
}

func (c *Client) Unmatch() {
	c.State = enums.NOT_IN_A_MATCH
	c.PairedAt = time.Time{}
	c.Pair = nil
}

func (c *Client) Match(client *Client) {
	c.State = enums.IN_A_MATCH
	c.PairedAt = time.Now()
	c.Pair = client
}
