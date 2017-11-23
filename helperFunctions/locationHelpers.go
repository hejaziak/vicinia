package helperFunctions

import (
	"errors"
	"strconv"
	"strings"

	datastructures "vicinia/datastructures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

//SetLocation : sets the location of client in locations map
func SetLocation(uuid uuid.UUID, location string) error {
	coordinates := strings.Split(location, ",")

	if err := datastructures.CreateLocationEntry(uuid, coordinates[0], coordinates[1]); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return errors.New("You have already entered your location successfully")
	}

	return nil
}

//CheckLocationFormat : checks that the location is entered in the correct format
func CheckLocationFormat(location string) error {
	if location == "" {
		return errors.New("empty location")
	}

	splittedLocation := strings.Split(location, ",")
	if len(splittedLocation) > 2 {
		return errors.New("more than one ',' in location")
	}

	if _, err := strconv.ParseFloat(splittedLocation[0], 64); err != nil {
		return errors.New("syntax error in latitude")
	}

	if _, err := strconv.ParseFloat(splittedLocation[1], 64); err != nil {
		return errors.New("syntax error in longitude")
	}

	return nil
}
