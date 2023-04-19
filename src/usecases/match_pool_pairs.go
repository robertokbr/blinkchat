package usecases

import (
	"time"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/utils"
)

type MatchPoolPairs struct {
	pool *models.Pool
}

func NewMatchPoolPairs(pool *models.Pool) *MatchPoolPairs {
	return &MatchPoolPairs{
		pool: pool,
	}
}

func (uc *MatchPoolPairs) genRandomIndexes(len int) (int, int) {
	i1 := utils.Rand(len)
	i2 := utils.Rand(len)

	if i1 == i2 {
		return uc.genRandomIndexes(len)
	}

	return i1, i2
}

func (uc *MatchPoolPairs) Execute() {
	for {
		amountOfPairs := len(uc.pool.Pairs)

		if amountOfPairs < 2 {
			// Wait for more clients for 5 seconds
			time.Sleep(5 * time.Second)
			continue
		}

		i1, i2 := uc.genRandomIndexes(amountOfPairs)

		c1 := uc.pool.Pairs[i1]
		c2 := uc.pool.Pairs[i2]

		logger.Infof("Matching clients %v and %v", c1.User.ID, c2.User.ID)

		c1.Match(c2)
		c2.Match(c1)

		utils.Splice(&uc.pool.Pairs, i1)
		utils.Splice(&uc.pool.Pairs, utils.If(i2 < i1, i2, i2-1))

		message := models.NewUserMatchedMessage(c2.User)

		if err := c1.Conn.WriteJSON(*message); err != nil {
			logger.Errorf("error writing message: %v", err)
		}

		message.Data.From = c1.User

		if err := c2.Conn.WriteJSON(*message); err != nil {
			logger.Errorf("error writing message: %v", err)
		}
	}
}
