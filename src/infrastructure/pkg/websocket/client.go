package websocket

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type Client struct {
	*models.User
	Conn  *websocket.Conn
	Pool  *Pool
	State enums.ClientState
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, contentInBytes, err := c.Conn.ReadMessage()

		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		content := string(contentInBytes)

		message := models.NewMessage(
			content,
			c.User,
			enums.TEXT,
			enums.BROADCASTING,
		)

		c.Pool.Broadcast <- *message
	}
}
