package helperFunctions

import (
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

//SimplifyList : returns the list of parameters which will be returned as a response message to the user's generic search
func SimplifyList(latitude string, longitude string, input []maps.PlacesSearchResult) ([]structures.Place, error) {
	output := make([]structures.Place, 5)

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

		output[i] = structures.Place{
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
