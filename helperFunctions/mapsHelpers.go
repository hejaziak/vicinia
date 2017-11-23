package helperFunctions

import (
	"context"

	globals "vicinia/globals"

	datastructures "vicinia/datastructures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
	"googlemaps.github.io/maps"
)

//GetList : returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func GetList(uuid uuid.UUID, keyword string) ([]maps.PlacesSearchResult, error) {
	//getting the Google Maps client
	mapsClient, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return nil, err
	}

	if keyword == "" {
		return nil, err
	}

	latitude, longitude, err := datastructures.GetLongLat(uuid)
	if err != nil {
		return nil, err
	}

	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: latitude, Lng: longitude},
		RankBy:   "distance",
		Keyword:  keyword,
	}

	res, err := mapsClient.NearbySearch(context.Background(), req)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		emptyList := make([]maps.PlacesSearchResult, 0)
		return emptyList, err
	}

	return res.Results, nil
}

//GetDetails : returns detailed information about a specific place
func GetDetails(uuid uuid.UUID, index int) (maps.PlaceDetailsResult, error) {
	c, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return maps.PlaceDetailsResult{}, err
	}

	placeID, err := datastructures.GetPlace(uuid, index)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return maps.PlaceDetailsResult{}, err
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	res, err := c.PlaceDetails(context.Background(), req)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return maps.PlaceDetailsResult{}, err
	}

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return maps.PlaceDetailsResult{}, err
	}

	return res, nil
}

//GetDistance : takes as an input the coordinates of the origin and the place id of the destination,
//and returns the distance in km between the origin and destination
func GetDistance(cord string, destination string) (string, error) {
	c, err := globals.GetMapClient()
	if err != nil {
		return "", err
	}

	req := &maps.DistanceMatrixRequest{
		Origins:      []string{cord},
		Destinations: []string{"place_id:" + destination},
		Units:        "metric",
	}

	res, err := c.DistanceMatrix(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.Rows[0].Elements[0].Distance.HumanReadable, nil
}
