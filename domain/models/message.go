package models

import (
	"time"

	"github.com/robertokbr/blinkchat/domain/enums"
)

type Message struct {
	ID        int               `json:"id"`
	Content   string            `json:"content"`
	Type      enums.MessageType `json:"type"`
	From      User              `json:"from"`
	CreatedAt time.Time         `json:"created_at"`
}
