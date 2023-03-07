package websocket

import (
	"log"
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan models.Message
	CreatedAt  time.Time
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan models.Message),
		CreatedAt:  time.Now(),
	}
}

func (pool *Pool) Start(poolNumber int) {
	log.Printf("[Pool %v]: Starting pool at %v", poolNumber, pool.CreatedAt)

	for {
		select {
		case client := <-pool.Register:
			log.Printf("[Pool %v]: Registering client %v", poolNumber, client.User.ID)

			pool.Clients[client] = true

			message := models.NewMessage(
				"New user joined",
				client.User,
				enums.TEXT,
				enums.CONNECTED,
			)

			go func() {
				for client := range pool.Clients {
					client.Conn.WriteJSON(message)
				}
			}()

			break
		case client := <-pool.Unregister:
			log.Printf("[Pool %v]: Unregistering client %v", poolNumber, client.User.ID)

			delete(pool.Clients, client)

			message := models.NewMessage(
				"User has disconnected",
				client.User,
				enums.TEXT,
				enums.DISCONNECTED,
			)

			go func() {
				for client := range pool.Clients {
					client.Conn.WriteJSON(message)
				}
			}()

			break
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Printf("[Pool %v]: error writing message: %v", poolNumber, err)
					return
				}
			}
		}
	}
}
