package models

import (
	"time"

	"github.com/robertokbr/blinkchat/domain/enums"
)

type Room struct {
	ID        int             `json:"id"`
	State     enums.RoomState `json:"state"`
	Users     []User          `json:"users"`
	Messages  []Message       `json:"messages"`
	CreatedAt time.Time       `json:"created_at"`
}
