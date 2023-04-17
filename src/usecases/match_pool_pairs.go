package usecases

import (
	"time"

	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/utils"
)

type MatchPoolPairs struct {
	Pool *models.Pool
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
		amountOfPairs := len(uc.Pool.Pairs)

		if amountOfPairs < 2 {
			// Wait for more clients for 5 seconds
			time.Sleep(5 * time.Second)
			continue
		}

		i1, i2 := uc.genRandomIndexes(amountOfPairs)

		c1 := uc.Pool.Pairs[i1]
		c2 := uc.Pool.Pairs[i2]

		logger.Infof("Matching clients %v and %v", c1.User.ID, c2.User.ID)

		c1.Match(c2)
		c2.Match(c1)

		utils.Splice(&uc.Pool.Pairs, i1)
		utils.Splice(&uc.Pool.Pairs, utils.If(i2 < i1, i2, i2-1))

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
