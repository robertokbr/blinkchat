package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type PoolWorker struct {
	pool *models.Pool
}

func NewPoolWorker(pool *models.Pool) *PoolWorker {
	return &PoolWorker{
		pool: pool,
	}
}

func (pw *PoolWorker) checkPair(client *models.Client) bool {
	return client.Pair != nil && pw.pool.Clients[client.Pair.ID] != nil
}

func (pw *PoolWorker) Execute(id int) {
	logger.Infof("[Pool %v]: Starting pool", id)

	for {
		select {
		case message := <-pw.pool.Broadcast:
			pair := pw.pool.Clients[message.Data.From.ID].Pair

			if pw.checkPair(pair) {
				if err := pair.Conn.WriteJSON(message); err != nil {
					logger.Errorf("[Pool %v]: error writing message: %v", id, err)
				}
			}

			break
		case message := <-pw.pool.Match:
			client := pw.pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving matching request from client %v", id, client.User.ID)

			if pw.checkPair(client) {
				client.Pair.Unmatch()

				userUnmatchedMessage := models.NewUserUnmatchedMessage(client.User)

				if err := client.Pair.Conn.WriteJSON(*userUnmatchedMessage); err != nil {
					logger.Errorf("error writing message: %v", err)
				}

				client.Unmatch()
			}

			pw.pool.Pairs = append(pw.pool.Pairs, client)

			client.State = enums.LOOKING_FOR_MATCH

			break
		case message := <-pw.pool.Unmatch:
			client := pw.pool.Clients[message.Data.From.ID]

			logger.Debugf("[Pool %v]: Receiving unmatching request from client %v", id, client.User.ID)

			if pw.checkPair(client) {
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
