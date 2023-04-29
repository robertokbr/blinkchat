package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

var RCWG = sync.WaitGroup{}

type RegisterClient struct {
	client *models.Client
	pool   *models.Pool
}

func NewRegisterClient(client *models.Client, pool *models.Pool) *RegisterClient {
	return &RegisterClient{
		client: client,
		pool:   pool,
	}
}

func (uc *RegisterClient) Execute() {
	logger.Infof("Registering client %v", uc.client.User.Email)

	uc.pool.Clients[uc.client.ID] = uc.client

	message := models.NewUserConnectedMessage(uc.client.User)

	RCWG.Add(1)

	go func() {
		defer RCWG.Done()

		for _, pc := range uc.pool.Clients {
			err := pc.Conn.WriteJSON(*message)

			if err != nil {
				logger.Errorf("error writing message: %v", err)
			}
		}
	}()
}
