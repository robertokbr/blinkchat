package websocket_tests

import (
	"testing"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
	websocket_test_factories "github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket/tests/factories"

	"github.com/stretchr/testify/require"
)

func init() {
	database.NewDatabase().Connect()
}

type WSConnectionSpy struct {
	MessagesSent         []models.Message
	WriteJSONTimesCalled int
}

func (wsc *WSConnectionSpy) WriteJSON(data interface{}) error {
	wsc.WriteJSONTimesCalled++
	wsc.MessagesSent = append(wsc.MessagesSent, *data.(*models.Message))
	return nil
}

func (wsc *WSConnectionSpy) Close() error {
	return nil
}

func (wsc *WSConnectionSpy) ReadMessage() (int, []byte, error) {
	return 0, nil, nil
}

func newWSConnectionSpy() *WSConnectionSpy {
	return &WSConnectionSpy{
		WriteJSONTimesCalled: 0,
	}
}

func TestRegisterClient(t *testing.T) {
	pool := websocket.NewPool()
	wsConnectionSpy := newWSConnectionSpy()
	users := websocket_test_factories.MakeTestUser(5)

	for _, user := range users {
		websocket.WG.Add(1)

		client := websocket.Client{
			Conn:  wsConnectionSpy,
			User:  user,
			Pool:  pool,
			State: enums.NOT_IN_A_MATCH,
		}

		pool.Register(&client)
		websocket.WG.Wait()
	}

	require.Equal(t, 5, len(pool.Clients))
	require.Equal(t, enums.CONNECTED, wsConnectionSpy.MessagesSent[0].Action)
	require.Equal(t, 15, wsConnectionSpy.WriteJSONTimesCalled)
}

func TestUnregisterClient(t *testing.T) {
	pool := websocket.NewPool()
	wsConnectionSpy := newWSConnectionSpy()
	users := websocket_test_factories.MakeTestUser(5)

	for _, user := range users {
		websocket.WG.Add(1)

		client := websocket.Client{
			Conn:  wsConnectionSpy,
			User:  user,
			Pool:  pool,
			State: enums.NOT_IN_A_MATCH,
		}

		pool.Register(&client)
		websocket.WG.Wait()
	}

	require.Equal(t, 5, len(pool.Clients))

	for _, client := range pool.Clients {
		websocket.WG.Add(1)
		pool.Unregister(client)
		websocket.WG.Wait()
	}

	require.Equal(t, 0, len(pool.Clients))
}

func TestUnregisterClientLookingForMatch(t *testing.T) {
	pool := websocket.NewPool()
	wsConnectionSpy := newWSConnectionSpy()
	users := websocket_test_factories.MakeTestUser(1)
	user := users[0]

	client := websocket.Client{
		Conn:  wsConnectionSpy,
		User:  user,
		Pool:  pool,
		State: enums.LOOKING_FOR_MATCH,
	}

	pool.Register(&client)
	pool.Pairs = append(pool.Pairs, &client)

	require.Equal(t, 1, len(pool.Clients))
	require.Equal(t, 1, len(pool.Pairs))

	pool.Unregister(&client)

	require.Equal(t, 0, len(pool.Clients))
	require.Equal(t, 0, len(pool.Pairs))
}

func TestUnregisterClientWithPair(t *testing.T) {
	pool := websocket.NewPool()
	wsConnectionSpy := newWSConnectionSpy()
	users := websocket_test_factories.MakeTestUser(2)
	userOne := users[0]
	userTwo := users[1]

	clientOne := websocket.Client{
		Conn:  wsConnectionSpy,
		User:  userOne,
		Pool:  pool,
		State: enums.NOT_IN_A_MATCH,
	}

	clientTwo := websocket.Client{
		Conn:  wsConnectionSpy,
		User:  userTwo,
		Pool:  pool,
		State: enums.NOT_IN_A_MATCH,
	}

	websocket.WG.Add(2)
	pool.Register(&clientOne)
	pool.Register(&clientTwo)
	websocket.WG.Wait()

	clientOne.Match(&clientTwo)
	clientTwo.Match(&clientOne)

	require.Equal(t, *clientOne.Pair, clientTwo)
	require.Equal(t, *clientTwo.Pair, clientOne)

	websocket.WG.Add(1)
	pool.Unregister(&clientOne)
	websocket.WG.Wait()

	userUnmatchedMessage := wsConnectionSpy.MessagesSent[len(wsConnectionSpy.MessagesSent)-2]

	require.Nil(t, clientOne.Pair)
	require.Nil(t, clientTwo.Pair)
	require.Equal(t, 1, len(pool.Clients))
	require.Equal(t, enums.NOT_IN_A_MATCH, clientTwo.State)
	require.Equal(t, enums.UNMATCHED, userUnmatchedMessage.Action)
}
