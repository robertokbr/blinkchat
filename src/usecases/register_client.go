package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

var RCWG = sync.WaitGroup{}

type RegisterClient struct {
	pool   *models.Pool
	client *models.Client
}

func NewRegisterClient(pool *models.Pool, client *models.Client) *RegisterClient {
	return &RegisterClient{
		pool:   pool,
		client: client,
	}
}

func (uc *RegisterClient) Execute() {
	logger.Infof("Registering client %v", uc.client.User.Email)

	uc.pool.Clients[uc.client.ID] = uc.client

	message := models.NewUserConnectedMessage(uc.client.User)

	go func() {
		defer RCWG.Done()

		for _, pc := range uc.pool.Clients {
			pc.Conn.WriteJSON(*message)
		}
	}()
}
