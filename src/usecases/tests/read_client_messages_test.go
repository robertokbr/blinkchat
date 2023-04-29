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

func simulateClientSendingMessages(message []byte, channel chan []byte) {
	channel <- message
	time.Sleep(1 * time.Second)
	simulateClientSendingMessages(message, channel)
}

func TestReadClientMessages(t *testing.T) {
	ws := usecases_tests_spies.NewWebsocketConnection()

	users := usecases_tests_factories.MakeTestUser(1)
	user := users[0]

	client := models.NewClient(user, ws)

	jsonMessage := "{ \"action\": \"broadcasting\",\"data\": { \"content\": \"Hello test\", \"message_type\": \"text\" } }"

	encodedMessage := []byte(jsonMessage)

	jobs := make(chan models.Message)

	readClientMessagesUsecase := usecases.NewReadClientMessages(client, jobs)

	go readClientMessagesUsecase.Execute()

	go simulateClientSendingMessages(encodedMessage, ws.Messages)

	poolMessage := <-jobs

	require.NotNil(t, poolMessage)
	require.Equal(t, "Hello test", poolMessage.Data.Content)
	require.Equal(t, enums.TEXT, poolMessage.Data.Type)
	require.Equal(t, enums.BROADCASTING, poolMessage.Action)
}
