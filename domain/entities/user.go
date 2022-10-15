package entities

import "github.com/robertokbr/blinkchat/domain/enums"

type User struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	ImageURL  string          `json:"image_url"`
	Email     string          `json:"email"`
	State     enums.UserState `json:"state"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}
