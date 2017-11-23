package globals

import (
	"errors"
	"strconv"
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
	if _, exists := locations[key]; exists == true {
		return errors.New("UUID already exists in locations map")
	}

	locations[key] = []string{latitude, longitude}

	time.AfterFunc(time.Duration(24*time.Hour), func() {
		if err := DeleteLocationEntry(key); err != nil {
			return
		}
	})

	return nil
}

//GetLocationEntry : returns an entry in locations map
func GetLocationEntry(key uuid.UUID) ([]string, error) {
	location, exists := locations[key]
	if exists == false {
		return nil, errors.New("UUID doesn't exist in locations map")
	}

	return location, nil
}

//GetLongLat : returns long and lat floats in locations map
func GetLongLat(key uuid.UUID) (float64, float64, error) {
	location, err := GetLocationEntry(key) //location contains latitude and longitude

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return 0, 0, errors.New("location not set")
	}

	latitude, err := strconv.ParseFloat(location[0], 64)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return 0, 0, err
	}

	longitude, err := strconv.ParseFloat(location[1], 64)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return 0, 0, err
	}

	return longitude, latitude, nil
}

//DeleteLocationEntry : deletes an entry in locations map
func DeleteLocationEntry(key uuid.UUID) error {
	if _, exists := locations[key]; exists == false {
		return errors.New("UUID doesn't exist in locations map")
	}

	delete(locations, key)
	return nil
}

//PrintLocationsMap : displays contents of locations map
func PrintLocationsMap() {
	pretty.Println("Cords -Map =========================")
	for key, value := range locations {
		pretty.Println("Key:", key, "Value:", value)
	}
	pretty.Println("====================================")
}
