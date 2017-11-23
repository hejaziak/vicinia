package globals

import (
	"errors"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client

//InitMapClient : intializes the MapsClient
func InitMapClient(apiKey string) error {
	mapsCandidate, err := maps.NewClient(maps.WithAPIKey(apiKey))

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return err
	}

	mapsClient = mapsCandidate
	return nil
}

//GetMapClient : returns the MapsClient
func GetMapClient() (*maps.Client, error) {
	if mapsClient == nil {
		return nil, errors.New("client not initialized")
	}

	return mapsClient, nil
}
