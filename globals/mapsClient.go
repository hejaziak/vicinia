package globals

import (
	"errors"
	"log"

	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client

//InitMapClient : intializes the MapsClient
func InitMapClient(apiKey string) {
	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

//GetMapClient : returns the MapsClient
func GetMapClient() (*maps.Client, error) {
	if mapsClient == nil {
		return nil, errors.New("client not initialized")
	}

	return mapsClient, nil
}
