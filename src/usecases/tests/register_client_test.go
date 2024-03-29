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

func TestRegisterClient(t *testing.T) {
	pool := models.NewPool()
	ws := usecases_tests_spies.NewWebsocketConnection()
	users := usecases_tests_factories.MakeTestUser(5)

	for _, user := range users {
		client := models.Client{
			Conn:  ws,
			User:  user,
			State: enums.NOT_IN_A_MATCH,
		}

		registerClientUsecase := usecases.NewRegisterClient(&client, pool)
		registerClientUsecase.Execute()
		usecases.RCWG.Wait()
	}

	require.Equal(t, 5, len(pool.Clients))
	require.Equal(t, enums.CONNECTED, ws.MessagesSent[0].Action)
	require.Equal(t, 15, ws.WriteJSONTimesCalled)
}
