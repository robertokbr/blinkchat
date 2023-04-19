package usecases_tests

import (
	"testing"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/usecases"
	usecases_tests_factories "github.com/robertokbr/blinkchat/src/usecases/tests/factories"
	usecases_tests_spies "github.com/robertokbr/blinkchat/src/usecases/tests/spies"

	"github.com/stretchr/testify/require"
)

func TestHandleEventMatching(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)
	matchPoolClientsUsecase := usecases.NewMatchPoolPairs(pool)
	u1 := users[0]
	u2 := users[1]

	c1 := models.Client{
		Conn:  ws,
		User:  u1,
		State: enums.LOOKING_FOR_MATCH,
	}

	c2 := models.Client{
		Conn:  ws,
		User:  u2,
		State: enums.LOOKING_FOR_MATCH,
	}

	pool.Clients[c1.ID] = &c1
	pool.Clients[c2.ID] = &c2

	pool.Pairs = append(pool.Pairs, &c1, &c2)

	go matchPoolClientsUsecase.Execute()

	for {
		if c1.State == enums.IN_A_MATCH && c2.State == enums.IN_A_MATCH {
			break
		}
	}

	require.Equal(t, *c1.Pair, c2)
	require.Equal(t, *c2.Pair, c1)
	require.Equal(t, enums.IN_A_MATCH, c1.State)
	require.Equal(t, enums.IN_A_MATCH, c2.State)
}
