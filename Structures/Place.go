package structures

type Place struct {
	Name         string  `json:"name"`
	Distance     string  `json:"distance"`
	Rating       float32 `json:"rating"`
	Type         string  `json:"type"`
	Address      string  `json:"address"`
	MobileNumber string  `json:"mobile_number"`
	Link         string  `json:"link"`
}

type PlaceListEntity struct {
	Name     string  `json:"name"`
	Distance string  `json:"distance"`
	Rating   float32 `json:"rating"`
}
