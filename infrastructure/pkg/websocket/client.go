package websocket

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/robertokbr/blinkchat/domain/models"
)

type Client struct {
	*models.User
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, content, err := c.Conn.ReadMessage()

		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		message := models.NewMessage(string(content), *c.User)

		c.Pool.Broadcast <- *message
	}
}
