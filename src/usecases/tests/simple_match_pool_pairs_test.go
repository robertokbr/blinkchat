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

func TestSimpleMatchPoolPairs(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(4)
	simpleMatchPoolPairsUsecase := usecases.NewSimpleMatchPoolPairs(pool)

	u1 := users[0]
	u2 := users[1]
	u3 := users[2]
	u4 := users[3]

	client1 := models.NewClient(u1, ws)
	client2 := models.NewClient(u2, ws)
	client3 := models.NewClient(u3, ws)
	client4 := models.NewClient(u4, ws)

	client1.State = enums.LOOKING_FOR_MATCH
	client2.State = enums.LOOKING_FOR_MATCH
	client3.State = enums.LOOKING_FOR_MATCH
	client4.State = enums.LOOKING_FOR_MATCH

	pool.Clients[client1.ID] = client1
	pool.Clients[client2.ID] = client2
	pool.Clients[client3.ID] = client3
	pool.Clients[client4.ID] = client4

	go simpleMatchPoolPairsUsecase.Execute()

	pool.Pairs <- client1
	pool.Pairs <- client2
	pool.Pairs <- client3

	for {
		if client1.State == enums.IN_A_MATCH && client2.State == enums.IN_A_MATCH {
			break
		}
	}

	pool.Pairs <- client4

	for {
		if client3.State == enums.IN_A_MATCH && client4.State == enums.IN_A_MATCH {
			break
		}
	}

	require.Equal(t, client1.Pair, client2)
	require.Equal(t, client2.Pair, client1)
	require.Equal(t, enums.IN_A_MATCH, client1.State)
	require.Equal(t, enums.IN_A_MATCH, client2.State)
}
