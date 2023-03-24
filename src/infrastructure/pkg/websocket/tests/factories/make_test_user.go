package websocker_test_factories

import (
	"fmt"

	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

func createUser(number int) *models.User {
	createConnectionDTO := dtos.CreateConnection{
		Name:     fmt.Sprintf("test%d", number),
		Email:    fmt.Sprintf("test%d@email.com", number),
		ImageURL: "https://test.com/image.png",
	}

	user := models.NewUser(createConnectionDTO)

	return user
}

func MakeTestUser(args ...interface{}) []*models.User {
	users := make([]*models.User, 0)

	if len(args) == 0 {
		user := createUser(1)

		users = append(users, user)
	} else {
		for i := 1; i <= args[0].(int); i++ {
			user := createUser(i)

			users = append(users, user)
		}
	}

	return users
}
