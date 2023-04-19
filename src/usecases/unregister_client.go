package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/utils"
)

var UCWG = sync.WaitGroup{}

type UnregisterClient struct {
	pool   *models.Pool
	client *models.Client
}

func NewUnregisterClient(pool *models.Pool, client *models.Client) *UnregisterClient {
	return &UnregisterClient{
		pool:   pool,
		client: client,
	}
}

func (uc *UnregisterClient) Execute() {
	logger.Infof("Unregistering client %v", uc.client.User.Email)

	// remove from pairs pool
	if uc.client.State == enums.LOOKING_FOR_MATCH {
		utils.Filter(&uc.pool.Pairs, func(c *models.Client) bool {
			return c.ID != uc.client.ID
		})
	}

	if uc.client.Pair != nil && uc.pool.Clients[uc.client.Pair.ID] != nil {
		userUnmatchedMessage := models.NewUserUnmatchedMessage(uc.client.User)

		uc.client.Pair.Unmatch()

		if err := uc.client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
			logger.Errorf("error writing message: %v", err)
		}
	}

	// remove from clients pool
	delete(uc.pool.Clients, uc.client.ID)

	message := models.NewUserDisconnectedMessage(uc.client.User)

	go func() {
		defer UCWG.Done()

		for _, pc := range uc.pool.Clients {
			pc.Conn.WriteJSON(*message)
		}
	}()
}
