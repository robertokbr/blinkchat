package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/messages"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/utils"
)

var WG = sync.WaitGroup{}

type Pool struct {
	Broadcast chan models.Message
	Match     chan models.Message
	Unmatch   chan models.Message
	Clients   map[string]*Client
	Pairs     []*Client
	CreatedAt time.Time
}

func NewPool() *Pool {
	return &Pool{
		Broadcast: make(chan models.Message),
		Match:     make(chan models.Message),
		Unmatch:   make(chan models.Message),
		Clients:   make(map[string]*Client),
		Pairs:     make([]*Client, 0),
		CreatedAt: time.Now(),
	}
}

func (pool *Pool) Register(client *Client) {
	logger.Infof("Registering client %v", client.User.Email)

	pool.Clients[client.ID] = client

	message := messages.UserConnected(client.User)

	go func() {
		defer WG.Done()

		for _, pc := range pool.Clients {
			pc.Conn.WriteJSON(message)
		}
	}()
}

func (pool *Pool) Unregister(client *Client) {
	logger.Infof("Unregistering client %v", client.User.Email)

	if client.State == enums.LOOKING_FOR_MATCH {
		utils.Filter(&pool.Pairs, func(c *Client) bool {
			return c.ID != client.ID
		})
	}

	pool.checkAndDisconnectPairs(client)
	delete(pool.Clients, client.ID)
	message := messages.UserDisconnected(client.User)

	go func() {
		defer WG.Done()

		for _, pc := range pool.Clients {
			pc.Conn.WriteJSON(message)
		}
	}()
}

func (pool *Pool) HandleEvent(message models.Message) {
	switch message.Action {
	case enums.BROADCASTING:
		pool.Broadcast <- message
		break
	case enums.MATCHING:
		pool.Match <- message
		break
	case enums.UNMATCHING:
		pool.Unmatch <- message
		break
	default:
		log.Printf("[debug]: Receiving message %v", message)
	}
}

func (pool *Pool) Start(poolNumber int) {
	logger.Infof("[Pool %v]: Starting pool", poolNumber)

	for {
		select {
		case message := <-pool.Broadcast:
			pair := pool.Clients[message.Data.From.ID].Pair

			if pair != nil && pool.checkIfClientIsOnline(pair) {
				if err := pair.Conn.WriteJSON(message); err != nil {
					logger.Errorf("[Pool %v]: error writing message: %v", poolNumber, err)
				}
			}

			break
		case message := <-pool.Match:
			client := pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving matching request from client %v", poolNumber, client.User.ID)

			pool.checkAndDisconnectPairs(client)

			pool.Pairs = append(pool.Pairs, client)

			client.State = enums.LOOKING_FOR_MATCH

			break
		case message := <-pool.Unmatch:
			client := pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving unmatching request from client %v", poolNumber, client.User.ID)

			pool.checkAndDisconnectPairs(client)
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

		logger.Infof("Matching clients %v and %v", clientOne.User.ID, clientTwo.User.ID)

		clientOne.Match(clientTwo)
		clientTwo.Match(clientOne)

		utils.Splice(&pool.Pairs, randomIndexOne)
		utils.Splice(&pool.Pairs, utils.If(randomIndexTwo < randomIndexOne, randomIndexTwo, randomIndexTwo-1))

		message := models.NewMessage(
			"You have been matched with a new user",
			clientTwo.User,
			enums.TEXT,
			enums.MATCHED,
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

func (pool *Pool) checkIfClientIsOnline(client *Client) bool {
	return pool.Clients[client.ID] != nil
}

func (pool *Pool) checkAndDisconnectPairs(client *Client) {
	if client.Pair != nil && pool.checkIfClientIsOnline(client.Pair) {
		userUnmatchedMessage := messages.UserUnmatched(client.User)
		client.Pair.Unmatch()
		if err := client.Pair.Conn.WriteJSON(userUnmatchedMessage); err != nil {
			logger.Errorf("error writing message: %v", err)
		}

		client.Unmatch()
	}
}
