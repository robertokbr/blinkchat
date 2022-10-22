package websocket

import (
	"log"
	"time"

	"github.com/robertokbr/blinkchat/domain/enums"
	"github.com/robertokbr/blinkchat/domain/models"
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

func (pool *Pool) Start() {
	log.Printf("Starting pool at %v", pool.CreatedAt)

	for {
		select {
		case client := <-pool.Register:
			log.Printf("Registering client %v", client.User.ID)

			pool.Clients[client] = true

			message := models.Message{
				Type:    enums.TEXT,
				Content: "New user joined",
			}

			for client := range pool.Clients {
				client.Conn.WriteJSON(message)
			}

			break

		case client := <-pool.Unregister:
			log.Printf("Unregistering client %v", client.User.ID)

			delete(pool.Clients, client)

			message := models.Message{
				Type:    enums.TEXT,
				Content: "User left",
			}

			for client := range pool.Clients {
				client.Conn.WriteJSON(message)
			}

			break
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Printf("error writing message: %v", err)
					return
				}
			}
		}
	}
}
