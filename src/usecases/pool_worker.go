package usecases

import (
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/utils"
)

type matchPoolPairsUsecase interface {
	Execute()
}

var maxMatchingRequestWaitInterval = 3

func checkPair(pool *models.Pool, client *models.Client) bool {
	return client.Pair != nil && pool.Clients[client.Pair.ID] != nil
}

func PoolWorker(id int, pool *models.Pool, jobs <-chan models.Message, matchPoolPairs matchPoolPairsUsecase) {
	logger.Infof("[Pool %v]: Starting pool", id)
	go matchPoolPairs.Execute()

	for message := range jobs {
		switch message.Action {
		case enums.BROADCASTING:
			pair := pool.Clients[message.Data.From.ID].Pair

			if checkPair(pool, pair) {
				if err := pair.Conn.WriteJSON(message); err != nil {
					logger.Errorf("[Pool %v]: error writing message: %v", id, err)
				}
			}

			break
		case enums.MATCHING:
			client := pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving matching request from client %v", id, client.User.ID)

			if checkPair(pool, client) {
				client.Pair.Unmatch()

				userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

				if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
					logger.Errorf("error writing message: %v", err)
				}

				client.Unmatch()
			}

			client.State = enums.LOOKING_FOR_MATCH

			randomInterval := time.Duration(utils.Rand(maxMatchingRequestWaitInterval)) * time.Second

			select {
			case <-time.After(randomInterval):
				pool.Pairs <- client
			}

			break
		case enums.UNMATCHING:
			client := pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving unmatching request from client %v", id, client.User.ID)

			if checkPair(pool, client) {
				client.Pair.Unmatch()

				userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

				if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
					logger.Errorf("error writing message: %v", err)
				}
			}

			client.Unmatch()
		}
	}

	logger.Infof("[Pool %v]: Closing pool", id)
}
