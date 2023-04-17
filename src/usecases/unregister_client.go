package usecases

import (
	"sync"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/utils"
)

type UnregisterClient struct {
	Pool *models.Pool
}

var UnregisterClientWG = sync.WaitGroup{}

func (uc *UnregisterClient) Execute(client *models.Client) {
	logger.Infof("Unregistering client %v", client.User.Email)

	// remove from pairs pool
	if client.State == enums.LOOKING_FOR_MATCH {
		utils.Filter(&uc.Pool.Pairs, func(c *models.Client) bool {
			return c.ID != client.ID
		})
	}

	if client.Pair != nil && uc.Pool.Clients[client.Pair.ID] != nil {
		userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

		client.Pair.Unmatch()

		if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
			logger.Errorf("error writing message: %v", err)
		}
	}

	// remove from clients pool
	delete(uc.Pool.Clients, client.ID)

	message := models.NewUserDisconnectedMessage(client.User)

	go func() {
		defer UnregisterClientWG.Done()

		for _, pc := range uc.Pool.Clients {
			pc.Conn.WriteJSON(*message)
		}
	}()
}
