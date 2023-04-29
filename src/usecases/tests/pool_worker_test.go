package usecases_tests

import (
	"testing"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/usecases"
	usecases_tests_factories "github.com/robertokbr/blinkchat/src/usecases/tests/factories"
	usecases_tests_fakes "github.com/robertokbr/blinkchat/src/usecases/tests/fakes"
	usecases_tests_spies "github.com/robertokbr/blinkchat/src/usecases/tests/spies"
	"github.com/stretchr/testify/require"
)

func TestPoolWorkerBroadcastCase(t *testing.T) {
	jobs := make(chan models.Message)
	pool := models.NewPool()

	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)

	u1 := users[0]
	u2 := users[1]

	c1 := models.NewClient(u1, ws)
	c2 := models.NewClient(u2, ws)

	c1.Match(c2)
	c2.Match(c1)

	pool.Clients[c1.ID] = c1
	pool.Clients[c2.ID] = c2

	message := models.NewPingMessage(u1)

	matchPoolPairs := usecases_tests_fakes.NewMatchPoolPairs()

	go usecases.PoolWorker(1, pool, jobs, matchPoolPairs)

	jobs <- *message

	for {
		if len(ws.MessagesSent) == 1 {
			break
		}
	}

	require.Equal(t, 1, len(ws.MessagesSent))
	require.Equal(t, message.Data.Content, ws.MessagesSent[0].Data.Content)
	require.Equal(t, message.Data.From, ws.MessagesSent[0].Data.From)
}

func TestPoolWorkerMatchCase(t *testing.T) {
	jobs := make(chan models.Message)
	pool := models.NewPool()

	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)

	u1 := users[0]

	c1 := models.NewClient(u1, ws)

	pool.Clients[c1.ID] = c1

	message := models.NewMessage(
		"MATCHING REQUEST",
		u1,
		enums.TEXT,
		enums.MATCHING,
	)

	matchPoolPairs := usecases_tests_fakes.NewMatchPoolPairs()

	go usecases.PoolWorker(1, pool, jobs, matchPoolPairs)

	jobs <- *message

	for {
		if c1.State == enums.LOOKING_FOR_MATCH {
			break
		}
	}

	waitingForMatchClient := <-pool.Pairs

	require.Equal(t, c1.State, enums.LOOKING_FOR_MATCH)
	require.Equal(t, waitingForMatchClient, c1)
}

func TestPoolWorkerMatchWithAlreadyMatchedClientCase(t *testing.T) {
	jobs := make(chan models.Message)
	pool := models.NewPool()

	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)

	u1 := users[0]
	u2 := users[1]

	c1 := models.NewClient(u1, ws)
	c2 := models.NewClient(u2, ws)

	c1.Match(c2)
	c2.Match(c1)

	pool.Clients[c2.ID] = c2
	pool.Clients[c1.ID] = c1

	message := models.NewMessage(
		"MATCHING REQUEST",
		u1,
		enums.TEXT,
		enums.MATCHING,
	)

	matchPoolPairs := usecases_tests_fakes.NewMatchPoolPairs()

	go usecases.PoolWorker(1, pool, jobs, matchPoolPairs)

	jobs <- *message

	for {
		if c1.State == enums.LOOKING_FOR_MATCH {
			break
		}
	}

	require.Equal(t, c1.State, enums.LOOKING_FOR_MATCH)
	require.Nil(t, c2.Pair)
	require.Equal(t, ws.MessagesSent[0].Action, enums.UNMATCHED)
}

func TestPoolWorkerUnmatchCase(t *testing.T) {
	jobs := make(chan models.Message)
	pool := models.NewPool()

	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)

	u1 := users[0]
	u2 := users[1]

	c1 := models.NewClient(u1, ws)

	c2 := models.NewClient(u2, ws)

	pool.Clients[c1.ID] = c1
	pool.Clients[c2.ID] = c2

	c1.Match(c2)
	c2.Match(c1)

	message := models.NewMessage(
		"UNMATCHING REQUEST",
		u1,
		enums.TEXT,
		enums.UNMATCHING,
	)

	matchPoolPairs := usecases_tests_fakes.NewMatchPoolPairs()

	go usecases.PoolWorker(1, pool, jobs, matchPoolPairs)

	jobs <- *message

	for {
		if len(ws.MessagesSent) == 1 {
			break
		}
	}

	require.Nil(t, c1.Pair)
	require.Nil(t, c2.Pair)
	require.Equal(t, c1.User, ws.MessagesSent[0].Data.From)
	require.Equal(t, enums.UNMATCHED, ws.MessagesSent[0].Action)
}
