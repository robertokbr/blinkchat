package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

var UCWG = sync.WaitGroup{}

type UnregisterClient struct {
	client *models.Client
	pool   *models.Pool
}

func NewUnregisterClient(client *models.Client, pool *models.Pool) *UnregisterClient {
	return &UnregisterClient{
		client: client,
		pool:   pool,
	}
}

func (uc *UnregisterClient) Execute() {
	logger.Infof("Unregistering client %v", uc.client.User.Email)

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

	UCWG.Add(1)

	go func() {
		defer UCWG.Done()

		for _, pc := range uc.pool.Clients {
			pc.Conn.WriteJSON(*message)
		}
	}()
}
