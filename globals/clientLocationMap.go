package globals

import (
	"errors"
	"time"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

var locations map[uuid.UUID][]string

//InitLocations : intializes the locations map
func InitLocations() {
	locations = make(map[uuid.UUID][]string)
}

//CreateLocationEntry : creates new entry in locations map
func CreateLocationEntry(key uuid.UUID, latitude string, longitude string) error {
	_, exists := locations[key]
	if(exists==true){
		return errors.New("UUID already exists in locations map")
	}

	locations[key] = []string{latitude,longitude}

	time.AfterFunc( time.Duration(24*time.Hour), func() {
		if err := DeleteLocationEntry(key); err != nil {
			return
		}
		
	})

	return nil
}

//GetLocationEntry : returns an entry in locations map
func GetLocationEntry(key uuid.UUID) ([]string, error) {
	location, exists := locations[key]
	if(exists==false){
		return nil, errors.New("UUID doesn't exist in locations map")
	}

	return location, nil
}

//DeleteLocationEntry : deletes an entry in locations map
func DeleteLocationEntry(key uuid.UUID) error {
	_, exists := locations[key]
	if(exists==false){
		return errors.New("UUID doesn't exist in locations map")
	}

	delete(locations, key)
	return nil
}

//PrintLocationsMap : displays contents of locations map
func PrintLocationsMap() {
	pretty.Println("Map ======================:")
	for key, value := range locations {
		pretty.Println("Key:", key, "Value:", value)
	}
	pretty.Println("====================================")
}