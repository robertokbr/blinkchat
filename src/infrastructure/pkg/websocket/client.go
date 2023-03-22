package websocket

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type Client struct {
	*models.User
	Conn     *websocket.Conn
	Pool     *Pool
	Pair     *Client
	PairedAt time.Time
}

func (c *Client) Unmatch() {
	c.PairedAt = time.Time{}
	c.Pair = nil
}

func (c *Client) Match(client *Client) {
	c.PairedAt = time.Now()
	c.Pair = client
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister(c)
		c.Conn.Close()
	}()

	for {
		_, websocketMessage, err := c.Conn.ReadMessage()

		if err != nil {
			logger.Errorf("error reading message: %v", err)
			break
		}

		createMessageDTO, err := dtos.NewCreateMessage(string(websocketMessage))

		if err != nil {
			continue
		}

		message := models.NewMessage(
			createMessageDTO.Content,
			c.User,
			enums.MessageType(createMessageDTO.MessageType),
			enums.WebsocketEvent(createMessageDTO.Event),
		)

		c.Pool.HandleEvent(*message)
	}
}
