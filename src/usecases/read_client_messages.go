package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type ReadClientMessages struct {
	client           *models.Client
	pool             *models.Pool
	unregisterClient *UnregisterClient
}

func NewReadClientMessages(
	pool *models.Pool,
	client *models.Client,
	unregisterClient *UnregisterClient,
) *ReadClientMessages {
	return &ReadClientMessages{
		client:           client,
		pool:             pool,
		unregisterClient: unregisterClient,
	}
}

func (uc *ReadClientMessages) Execute() {
	defer func() {
		uc.unregisterClient.Execute()
		uc.client.Conn.Close()
	}()

	for {
		_, websocketMessage, err := uc.client.Conn.ReadMessage()

		if err != nil {
			logger.Errorf("error reading message: %v", err)
			break
		}

		createMessageDTO, err := dtos.NewCreateMessage(string(websocketMessage))

		if err != nil {
			logger.Errorf("error parsing message: %v", err)
			continue
		}

		message := models.NewMessage(
			createMessageDTO.Content,
			uc.client.User,
			enums.MessageType(createMessageDTO.MessageType),
			enums.WebsocketEvent(createMessageDTO.Event),
		)

		uc.pool.PushMessage(*message)
	}
}
