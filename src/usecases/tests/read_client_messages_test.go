package usecases_tests

import (
	"testing"
	"time"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/usecases"
	usecases_tests_factories "github.com/robertokbr/blinkchat/src/usecases/tests/factories"
	usecases_tests_spies "github.com/robertokbr/blinkchat/src/usecases/tests/spies"

	"github.com/stretchr/testify/require"
)

func simulateClientSendingMessages(message string, channel chan string) {
	channel <- message
	time.Sleep(1 * time.Second)
	simulateClientSendingMessages(message, channel)
}

func TestReadClientMessages(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(1)
	user := users[0]

	client := models.Client{
		Conn:  ws,
		User:  user,
		State: enums.NOT_IN_A_MATCH,
	}

	jsonMessage := "{ \"action\": \"broadcasting\",\"data\": { \"content\": \"Hello test\", \"message_type\": \"text\" } }"

	go simulateClientSendingMessages(jsonMessage, ws.Messages)

	readClientMessagesUsecase := usecases.ReadClientMessages{
		Client:           &client,
		Pool:             pool,
		UnregisterClient: &usecases.UnregisterClient{Pool: pool},
	}

	go readClientMessagesUsecase.Execute()

	poolMessage := <-pool.Broadcast

	require.NotNil(t, poolMessage)
	require.Equal(t, "Hello test", poolMessage.Data.Content)
	require.Equal(t, enums.TEXT, poolMessage.Data.Type)
	require.Equal(t, enums.BROADCASTING, poolMessage.Action)
}
