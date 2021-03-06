package helperFunctions

import (
	"context"
	"errors"

	globals "vicinia/globals"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

//GetList : returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func GetList(latitude float64, longitude float64, keyword string) ([]maps.PlacesSearchResult, error) {
	if keyword == "" {
		return nil, errors.New("empty keywords")
	}

	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: latitude, Lng: longitude},
		RankBy:   "distance",
		Keyword:  keyword,
	}

	mapsClient, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return nil, err
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
func GetDetails(placeID string) (maps.PlaceDetailsResult, error) {
	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	mapsClient, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return maps.PlaceDetailsResult{}, err
	}

	res, err := mapsClient.PlaceDetails(context.Background(), req)
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
	req := &maps.DistanceMatrixRequest{
		Origins:      []string{cord},
		Destinations: []string{"place_id:" + destination},
		Units:        "metric",
	}

	mapsClient, err := globals.GetMapClient()
	if err != nil {
		return "", err
	}

	res, err := mapsClient.DistanceMatrix(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.Rows[0].Elements[0].Distance.HumanReadable, nil
}
