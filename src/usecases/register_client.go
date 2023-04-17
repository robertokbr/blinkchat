package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type RegisterClient struct {
	Pool *models.Pool
}

var RegisterClientWG = sync.WaitGroup{}

func (uc *RegisterClient) Execute(client *models.Client) {
	logger.Infof("Registering client %v", client.User.Email)

	uc.Pool.Clients[client.ID] = client

	message := models.NewUserConnectedMessage(client.User)

	go func() {
		defer RegisterClientWG.Done()

		for _, pc := range uc.Pool.Clients {
			pc.Conn.WriteJSON(*message)
		}
	}()
}
