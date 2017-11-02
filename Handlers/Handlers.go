package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	structures "Vicinia/Structures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"github.com/kamalpy/apiai-go"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {

	welcomeMessage := structures.WelcomeStruct{
		Message: "Welcome ,where do you want to go ?",
		UUID:    uuid.NewV1(),
	}

	if err := json.NewEncoder(w).Encode(welcomeMessage); err != nil {
		panic(err)
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	ai := apiaigo.APIAI{
		AuthToken: "71027bbaf70a4a53847bedce6b83c94f",
		Language:  "en-US",
		SessionID: "1234567890",
	}

	resp, err := ai.SendText("I want restraunts")

	keyword := resp.Result.Parameters["keyword"]
	fmt.Println(keyword);


	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: 29.985352, Lng: 31.279194},
		RankBy:   "distance",
		Keyword:  keyword,
	}

	res, err := c.NearbySearch(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	output := SimplifyList(res.Results)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		panic(err)
	}
}

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID: "ChIJZ711I304WBQRZi-S-IYgfTE",
	}

	res, err := c.PlaceDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	output := SimplifyDetails(res)
	pretty.Println(res)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		panic(err)
	}
}

func SimplifyList(input []maps.PlacesSearchResult) []structures.PlaceListEntity {
	output := make([]structures.PlaceListEntity, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			break
		}

		output[i] = structures.PlaceListEntity{
			Name:     input[i].Name,
			Distance: getDistance("29.985352,31.279194", input[i].PlaceID),
			Rating:   input[i].Rating,
		}
	}

	return output
}

func SimplifyDetails(input maps.PlaceDetailsResult) structures.Place {
	output := structures.Place{
		Name:         input.Name,
		Distance:     getDistance("29.985352,31.279194", input.PlaceID),
		Rating:       input.Rating,
		Type:         input.Types[0],
		Address:      input.FormattedAddress,
		MobileNumber: input.FormattedPhoneNumber,
		Link:         input.URL,
	}
	return output
}

func getDistance(cord string, destination string) string {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	req := &maps.DistanceMatrixRequest{
		Origins:      []string{cord},
		Destinations: []string{"place_id:" + destination},
		Units:        "metric",
	}

	res, err := c.DistanceMatrix(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return res.Rows[0].Elements[0].Distance.HumanReadable
}
