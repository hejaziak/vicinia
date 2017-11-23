package structures

import (
	"github.com/satori/go.uuid"
)

type UUIDMessage struct {
	Message string    `json:"message"`
	UUID    uuid.UUID `json:"uuid"`
}

type Message struct {
	Message string `json:"message"`
}
