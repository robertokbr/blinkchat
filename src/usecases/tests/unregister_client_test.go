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

func TestUnregisterClient(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(5)

	for _, user := range users {
		client := models.Client{
			Conn:  ws,
			User:  user,
			State: enums.NOT_IN_A_MATCH,
		}

		usecases.RCWG.Add(1)
		registerClientUsecase := usecases.NewRegisterClient(pool, &client)
		registerClientUsecase.Execute()
		usecases.RCWG.Wait()
	}

	require.Equal(t, 5, len(pool.Clients))

	for _, client := range pool.Clients {
		usecases.UCWG.Add(1)
		unregisterClientUsecase := usecases.NewUnregisterClient(pool, client)
		unregisterClientUsecase.Execute()
		usecases.UCWG.Wait()
	}

	lastMessage := ws.MessagesSent[ws.WriteJSONTimesCalled-1].Action

	require.Equal(t, 0, len(pool.Clients))
	require.Equal(t, enums.DISCONNECTED, lastMessage)
	require.Equal(t, 25, ws.WriteJSONTimesCalled)
}

func TestUnregisterClientLookingForMatch(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(1)
	user := users[0]

	client := models.Client{
		Conn:  ws,
		User:  user,
		State: enums.LOOKING_FOR_MATCH,
	}

	unregisterClientUsecase := usecases.NewUnregisterClient(pool, &client)
	pool.Clients[user.ID] = &client
	pool.Pairs = append(pool.Pairs, &client)

	require.Equal(t, 1, len(pool.Clients))
	require.Equal(t, 1, len(pool.Pairs))

	usecases.UCWG.Add(1)
	unregisterClientUsecase.Execute()
	usecases.UCWG.Wait()

	require.Equal(t, 0, len(pool.Clients))
	require.Equal(t, 0, len(pool.Pairs))
}

func TestUnregisterClientWithPair(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(2)
	u1 := users[0]
	u2 := users[1]

	c1 := models.Client{
		Conn:  ws,
		User:  u1,
		State: enums.NOT_IN_A_MATCH,
	}

	c2 := models.Client{
		Conn:  ws,
		User:  u2,
		State: enums.NOT_IN_A_MATCH,
	}

	unregisterClientUsecase := usecases.NewUnregisterClient(pool, &c1)
	pool.Clients[u1.ID] = &c1
	pool.Clients[u2.ID] = &c2
	c1.Match(&c2)
	c2.Match(&c1)

	require.Equal(t, *c1.Pair, c2)
	require.Equal(t, *c2.Pair, c1)

	usecases.UCWG.Add(1)
	unregisterClientUsecase.Execute()
	usecases.UCWG.Wait()

	lastMessageSent := ws.MessagesSent[len(ws.MessagesSent)-2]

	require.Nil(t, c2.Pair)
	require.Equal(t, 1, len(pool.Clients))
	require.Equal(t, enums.NOT_IN_A_MATCH, c2.State)
	require.Equal(t, enums.UNMATCHED, lastMessageSent.Action)
}