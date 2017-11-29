package helperFunctions

import (
	"strconv"

	structures "vicinia/structures"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

//SimplifyList : returns the list of parameters which will be returned as a response message to the user's generic search
func SimplifyList(latitude string, longitude string, input []maps.PlacesSearchResult) ([]structures.PlaceListEntity, error) {
	output := make([]structures.PlaceListEntity, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			output = output[:i]
			break
		}

		name := input[i].Name
		if name == "" {
			name = "not specified"
		}

		distance, err := GetDistance(latitude+","+longitude, input[i].PlaceID)
		if err != nil {
			return nil, err
		}

		output[i] = structures.PlaceListEntity{
			Name:     name,
			Distance: distance,
			Rating:   input[i].Rating,
			ID:       input[i].PlaceID,
		}
	}

	return output, nil
}

//SimplifyDetails : returns the list of parameters which will be returned as a response message to the user's specific query
func SimplifyDetails(latitude string, longitude string, input maps.PlaceDetailsResult) (structures.Place, error) {
	name := input.Name
	if name == "" {
		name = "not specified"
	}

	distance, err := GetDistance(latitude+","+longitude, input.PlaceID)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return structures.Place{}, err
	}

	types := input.Types[0]
	if types == "" {
		types = "not specified"
	}

	address := input.FormattedAddress
	if address == "" {
		address = "not specified"
	}

	phone := input.FormattedPhoneNumber
	if phone == "" {
		phone = "not specified"
	}

	url := input.URL
	if url == "" {
		url = "not specified"
	}

	output := structures.Place{
		Name:         name,
		Distance:     distance,
		Rating:       input.Rating,
		Type:         types,
		Address:      address,
		MobileNumber: phone,
		Link:         url,
	}
	return output, nil
}

/*
//formatList: returns a formatted message containing list of places
func formatList(placesList []structures.PlaceListEntity, message string) structures.Message {
	formattedMessage := ""
	for _, place := range placesList {
		formattedMessage +=
			"Name: " + place.Name + " <br/> " +
				"Distance: " + place.Distance + " <br/> " +
				"Rating: " + strconv.FormatFloat(float64(place.Rating), 'f', -1, 32) + " <br/> " +
				"ID: " + strconv.Itoa(place.ID) + " <br/> <br/> "
	}

	formattedMessage += message

	return structures.Message{
		Message: formattedMessage,
	}
}
*/

//formatDetails: returns a formatted message containing details of a specifice place
func formatDetails(placeDetails structures.Place, message string) structures.Message {
	formattedMessage :=
		"Name: " + placeDetails.Name + " <br/> " +
			"Distance: " + placeDetails.Distance + " <br/> " +
			"Rating: " + strconv.FormatFloat(float64(placeDetails.Rating), 'f', -1, 32) + " <br/> " +
			"Type: " + placeDetails.Type + " <br/> " +
			"Address: " + placeDetails.Address + " <br/> " +
			"MobileNumber: " + placeDetails.MobileNumber + " <br/> " +
			"Link: <a href= " + placeDetails.Link + " > google maps </a>" + " <br/> <br/> " +
			message

	return structures.Message{
		Message: formattedMessage,
	}
}
