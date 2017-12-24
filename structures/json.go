package structures

import (
	"github.com/satori/go.uuid"
)

//Message : json message
type Message struct {
	Message string `json:"message"`
}

//UUIDMessage : json message including a uuid
type UUIDMessage struct {
	UUID    uuid.UUID `json:"uuid"`
	Message string    `json:"message"`
}

//LatLongMessage : json message including location
type LatLongMessage struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Message   string `json:"message"`
}

//PlaceListMessage : json message including an array of lists
type PlaceListMessage struct {
	Message []Place `json:"list"`
}
