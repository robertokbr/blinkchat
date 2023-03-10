package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type Client struct {
	*models.User
	Conn     *websocket.Conn
	Pool     *Pool
	State    enums.ClientState
	Pair     *Client
	PairedAt time.Time
}

func (c *Client) Unmatch() {
	c.State = enums.UNMATCHED
	c.PairedAt = time.Time{}
	c.Pair = nil
}

func (c *Client) Match(client *Client) {
	c.State = enums.MATCHED
	c.PairedAt = time.Now()
	c.Pair = client
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, websocketMessage, err := c.Conn.ReadMessage()

		if err != nil {
			log.Printf("error reading message: %v", err)
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

		c.Pool.SpreadMessage(*message)
	}
}
