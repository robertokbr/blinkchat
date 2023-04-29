package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type SimpleMatchPoolPairs struct {
	pool *models.Pool
}

func NewSimpleMatchPoolPairs(pool *models.Pool) *SimpleMatchPoolPairs {
	return &SimpleMatchPoolPairs{
		pool: pool,
	}
}

func (uc *SimpleMatchPoolPairs) checkMatch(client *models.Client, match *models.Client) bool {
	isClientOK := uc.pool.Clients[client.ID] != nil || client.State == enums.LOOKING_FOR_MATCH
	isMatchOK := uc.pool.Clients[match.ID] != nil || match.State == enums.LOOKING_FOR_MATCH

	if isClientOK && !isMatchOK {
		uc.pool.Pairs <- client
		return false
	}

	if isMatchOK && !isClientOK {
		uc.pool.Pairs <- match
		return false
	}

	if !isClientOK && !isMatchOK {
		return false
	}

	return true
}

func (uc *SimpleMatchPoolPairs) getMatch(client *models.Client) *models.Client {
	match := <-uc.pool.Pairs

	if client.ID == match.ID {
		return uc.getMatch(client)
	}

	return match
}

func (uc *SimpleMatchPoolPairs) Execute() {
	for client := range uc.pool.Pairs {
		match := uc.getMatch(client)

		if ok := uc.checkMatch(client, match); !ok {
			continue
		}

		logger.Infof("Matching clients %v and %v", client.User.ID, match.User.ID)

		client.Match(match)
		match.Match(client)

		message := models.NewUserMatchedMessage(match.User)

		if err := client.Conn.WriteJSON(*message); err != nil {
			logger.Errorf("error writing message: %v", err)
		}

		message.Data.From = client.User

		if err := match.Conn.WriteJSON(*message); err != nil {
			logger.Errorf("error writing message: %v", err)
		}
	}
}
