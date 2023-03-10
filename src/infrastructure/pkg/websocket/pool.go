package websocket

import (
	"log"
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/utils"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Pairs      []*Client
	Broadcast  chan models.Message
	Match      chan models.Message
	CreatedAt  time.Time
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Pairs:      make([]*Client, 0),
		Broadcast:  make(chan models.Message),
		Match:      make(chan models.Message),
		CreatedAt:  time.Now(),
	}
}

func (pool *Pool) SpreadMessage(message models.Message) {
	switch message.Action {
	case enums.BROADCASTING:
		pool.Broadcast <- message
		break
	case enums.MATCHING:
		pool.Match <- message
		break
	default:
		log.Printf("[debug]: Receiving message %v", message)
	}
}

func (pool *Pool) CheckIfClientIsOnline(client *Client) bool {
	return pool.Clients[client.ID] != nil
}

func (pool *Pool) Start(poolNumber int) {
	log.Printf("[Pool %v]: Starting pool at %v", poolNumber, pool.CreatedAt)

	for {
		select {
		case client := <-pool.Register:
			log.Printf("[Pool %v]: Registering client %v", poolNumber, client.User.ID)

			pool.Clients[client.ID] = client

			message := models.NewMessage(
				"New user joined",
				client.User,
				enums.TEXT,
				enums.CONNECTED,
			)

			go func() {
				for _, pc := range pool.Clients {
					pc.Conn.WriteJSON(message)
				}
			}()

			break
		case client := <-pool.Unregister:
			log.Printf("[Pool %v]: Unregistering client %v", poolNumber, client.User.ID)

			delete(pool.Clients, client.ID)

			// Improve this logic performance
			utils.Filter(&pool.Pairs, func(c *Client) bool {
				return c.ID != client.ID
			})

			message := models.NewMessage(
				"User has disconnected",
				client.User,
				enums.TEXT,
				enums.DISCONNECTED,
			)

			go func() {
				if client.Pair != nil && pool.CheckIfClientIsOnline(client.Pair) {
					// Unmatch the pair
					client.Pair.Unmatch()
					pool.Pairs = append(pool.Pairs, client.Pair)
				}

				for _, pc := range pool.Clients {
					pc.Conn.WriteJSON(message)
				}
			}()

			break
		case message := <-pool.Broadcast:
			pair := pool.Clients[message.Data.From.ID].Pair

			if pair != nil && pool.CheckIfClientIsOnline(pair) {
				if err := pair.Conn.WriteJSON(message); err != nil {
					log.Printf("[Pool %v]: error writing message: %v", poolNumber, err)
				}
			}

			break
		case message := <-pool.Match:
			client := pool.Clients[message.Data.From.ID]

			log.Printf("[debug]: Receiving matching request from client %v", client.User.ID)

			pair := client.Pair

			if pair != nil && pool.CheckIfClientIsOnline(pair) {
				// Notify pair that the user is searching for a new pair
				if err := pair.Conn.WriteJSON(message); err != nil {
					log.Printf("[Pool %v]: error writing message: %v", poolNumber, err)
				}

				pair.Unmatch()
			}

			client.Unmatch()

			pool.Pairs = append(pool.Pairs, client)

			break
		}
	}
}

func (pool *Pool) MatchPairs() {
	for {
		amountOfPairs := len(pool.Pairs)

		if amountOfPairs < 2 {
			// Wait for more clients for 5 seconds
			time.Sleep(5 * time.Second)
			continue
		}

		randomIndexOne, randomIndexTwo := pool.getTwoRandomIndex(amountOfPairs)

		clientOne := pool.Pairs[randomIndexOne]
		clientTwo := pool.Pairs[randomIndexTwo]

		log.Printf("Matching clients %v and %v", clientOne.User.ID, clientTwo.User.ID)

		clientOne.Match(clientTwo)
		clientTwo.Match(clientOne)

		utils.Splice(&pool.Pairs, randomIndexOne)
		utils.Splice(&pool.Pairs, utils.If(randomIndexTwo < randomIndexOne, randomIndexTwo, randomIndexTwo-1))

		message := models.NewMessage(
			"You have been matched with a new user",
			clientTwo.User,
			enums.TEXT,
			enums.MATCHING,
		)

		if err := clientOne.Conn.WriteJSON(message); err != nil {
			log.Printf("[error]: error writing message: %v", err)
		}

		message.Data.From = clientOne.User

		if err := clientTwo.Conn.WriteJSON(message); err != nil {
			log.Printf("[error]: error writing message: %v", err)
		}
	}
}

func (pool *Pool) getTwoRandomIndex(len int) (int, int) {
	randomIndexOne := utils.Rand(len)
	randomIndexTwo := utils.Rand(len)

	if randomIndexOne == randomIndexTwo {
		return pool.getTwoRandomIndex(len)
	}

	return randomIndexOne, randomIndexTwo
}
