package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type ReadClientMessages struct {
	Client           *models.Client
	Pool             *models.Pool
	UnregisterClient *UnregisterClient
}

func (uc *ReadClientMessages) Execute() {
	defer func() {
		uc.UnregisterClient.Execute(uc.Client)
		uc.Client.Conn.Close()
	}()

	for {
		_, websocketMessage, err := uc.Client.Conn.ReadMessage()

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
			uc.Client.User,
			enums.MessageType(createMessageDTO.MessageType),
			enums.WebsocketEvent(createMessageDTO.Event),
		)

		uc.Pool.PushMessage(*message)
	}
}
