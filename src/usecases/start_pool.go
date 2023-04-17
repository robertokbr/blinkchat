package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type StartPool struct {
	Pool *models.Pool
}

func (ps *StartPool) checkPair(client *models.Client) bool {
	return client.Pair != nil && ps.Pool.Clients[client.Pair.ID] != nil
}

func (ps *StartPool) Execute(poolNumber int) {
	logger.Infof("[Pool %v]: Starting pool", poolNumber)

	for {
		select {
		case message := <-ps.Pool.Broadcast:
			pair := ps.Pool.Clients[message.Data.From.ID].Pair

			if ps.checkPair(pair) {
				if err := pair.Conn.WriteJSON(message); err != nil {
					logger.Errorf("[Pool %v]: error writing message: %v", poolNumber, err)
				}
			}

			break
		case message := <-ps.Pool.Match:
			client := ps.Pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving matching request from client %v", poolNumber, client.User.ID)

			if ps.checkPair(client) {
				client.Pair.Unmatch()

				userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

				if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
					logger.Errorf("error writing message: %v", err)
				}

				client.Unmatch()
			}

			ps.Pool.Pairs = append(ps.Pool.Pairs, client)

			client.State = enums.LOOKING_FOR_MATCH

			break
		case message := <-ps.Pool.Unmatch:
			client := ps.Pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving unmatching request from client %v", poolNumber, client.User.ID)

			if ps.checkPair(client) {
				client.Pair.Unmatch()

				userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

				if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
					logger.Errorf("error writing message: %v", err)
				}
			}

			client.Unmatch()
		}
	}
}
