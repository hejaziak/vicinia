package structures

type Place struct {
	Name         string  `json:"name,omitempty"`
	Distance     string  `json:"distance,omitempty"`
	Rating       float32 `json:"rating,omitempty"`
	ID           string  `json:"id,omitempty"`
	Type         string  `json:"type,omitempty"`
	Address      string  `json:"address,omitempty"`
	MobileNumber string  `json:"mobile_number,omitempty"`
	Link         string  `json:"link,omitempty"`
}
