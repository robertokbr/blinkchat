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

type WSConnectionSpy struct {
	LastMessageSent      models.Message
	WriteJSONTimesCalled int
}

func (wsc *WSConnectionSpy) WriteJSON(data interface{}) error {
	wsc.WriteJSONTimesCalled++
	wsc.LastMessageSent = *data.(*models.Message)
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

func init() {
	database.NewDatabase().Connect()
}

func TestRegisterClient(t *testing.T) {
	pool := websocket.NewPool()

	users := websocket_test_factories.MakeTestUser(5)

	wsConnectionSpy := newWSConnectionSpy()

	channel := make(chan int, 0)

	for i, user := range users {
		client := websocket.Client{
			Conn:  wsConnectionSpy,
			User:  user,
			Pool:  pool,
			State: enums.NOT_IN_A_MATCH,
		}

		pool.Register(&client, channel, i)
	}

	require.Equal(t, 5, len(pool.Clients))
	require.NotNil(t, wsConnectionSpy.LastMessageSent.Action)
	require.Equal(t, enums.CONNECTED, wsConnectionSpy.LastMessageSent.Action)
	require.Equal(t, 20, wsConnectionSpy.WriteJSONTimesCalled)
}
