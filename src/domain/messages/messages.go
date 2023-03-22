package messages

import (
	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

func UserDisconnected(user *models.User) *models.Message {
	return models.NewMessage(
		"User has disconnected",
		user,
		enums.TEXT,
		enums.DISCONNECTED,
	)
}

func UserConnected(user *models.User) *models.Message {
	return models.NewMessage(
		"New user joined",
		user,
		enums.TEXT,
		enums.CONNECTED,
	)
}
