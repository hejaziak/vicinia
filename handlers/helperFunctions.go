package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	globals "vicinia/globals"
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/marcossegovia/apiai-go"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

//getList: returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func getList(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, message string) {
	c, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	client, err := apiai.NewClient(
		&apiai.ClientConfig{
			Token:      "71027bbaf70a4a53847bedce6b83c94f",
			QueryLang:  "en",    //Default en
			SpeechLang: "en-US", //Default en-US
		},
	)

	if err != nil {
		returnError(w, "")
	}

	//Set the query string and your current user identifier.
	qr, err := client.Query(apiai.Query{Query: []string{message}, SessionId: uuid.String()})
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	action := qr.Result.Action

	if strings.Compare(action, "search") != 0 {
		res := structures.Message{
			Message: qr.Result.Fulfillment.Speech,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

	} else {

		keyword := qr.Result.Params["keyword"]

		req := &maps.NearbySearchRequest{
			Location: &maps.LatLng{Lat: 29.985352, Lng: 31.279194},
			RankBy:   "distance",
			Keyword:  string(keyword.(string)),
		}

		res, err := c.NearbySearch(context.Background(), req)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "sorry, I couldn't find any relevant places")
			return
		}
		output, err := SimplifyList(res.Results)

		jsonMessage, _ := json.Marshal(output)

		respondMessage := extractMessage(string(jsonMessage), "To get detailed information about a specific place, please type its ID")
		if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

		inUUID, err := extractUUID(r)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			return
		}

		updateSession(inUUID, res.Results)
	}
}

//getDetails: returns detailed information about a specific place
func getDetails(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, index int) {
	c, err := globals.GetMapClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	placeID, err := globals.GetPlace(uuid, index)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	res, err := c.PlaceDetails(context.Background(), req)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	output, err := SimplifyDetails(res)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}
	pretty.Println(res)

	jsonMessage, _ := json.Marshal(output)
	respondMessage := extractMessage(string(jsonMessage), "Any other place you want to search for ?")

	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}
}

//extractUUID: returns uuid wich is extraced from the header of the request
func extractUUID(r *http.Request) (uuid.UUID, error) {

	s := r.Header.Get("Authorization")
	if s == "" {
		pretty.Printf("fatal error: Authorization Header empty \n")
		return uuid.Nil, errors.New("Authorization Header empty")
	}

	inUUID, err := uuid.FromString(s)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return uuid.Nil, err
	}

	if _, err := globals.GetEntry(inUUID); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return uuid.Nil, errors.New("sorry, UUID is not correct, please access /welcome to receive an UUID")
	}

	return inUUID, nil
}

//updateSession: updates the last places returned to the user
func updateSession(UUID uuid.UUID, input []maps.PlacesSearchResult) error {
	placeIDs := make([]string, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			break
		}
		placeIDs[i] = input[i].PlaceID
	}

	err := globals.UpdateEntry(UUID, placeIDs)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return err
	}

	return nil
}

func returnUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)

	if message == "" {
		message = "Oops! something went wrong"
	}

	respondMessage := structures.Message{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
	}
}

//returnError: returns error message
func returnError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Oops! something went wrong"
	}

	respondMessage := structures.Message{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
	}
}

//SimplifyList : returns the list of parameters which will be returned as a response message to the user's generic search
func SimplifyList(input []maps.PlacesSearchResult) ([]structures.PlaceListEntity, error) {
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

		distance, err := getDistance("29.985352,31.279194", input[i].PlaceID) //client coordinates are hard-coded for now
		if err != nil {
			return nil, err
		}

		output[i] = structures.PlaceListEntity{
			Name:     name,
			Distance: distance,
			Rating:   input[i].Rating,
			ID:       i + 1,
		}
	}

	return output, nil
}

//SimplifyDetails: returns the list of parameters which will be returned as a response message to the user's specific query
func SimplifyDetails(input maps.PlaceDetailsResult) (structures.Place, error) {
	name := input.Name
	if name == "" {
		name = "not specified"
	}

	distance, err := getDistance("29.985352,31.279194", input.PlaceID) //client coordinates are hard-coded for now
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

//getDistance: takes as an input the coordinates of the origin and the place id of the destination,
//and returns the distance in km between the origin and destination
func getDistance(cord string, destination string) (string, error) {
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

//extractMessage: returns a formatted response message to be readable by the client
func extractMessage(json string, message string) structures.Message {
	s2 := strings.Replace(json, "{", "", -1)
	s3 := strings.Replace(s2, "}", "", -1)
	s4 := strings.Replace(s3, "[", "", -1)
	s5 := strings.Replace(s4, "]", "", -1)
	s6 := strings.Replace(s5, "\"", "", -1)
	s7 := strings.Replace(s6, ",", " <br/> ", -1)
	cleanString := strings.Replace(s7, "<br/> name:", "<br/> <br/> name:", -1)
	cleanString = cleanString + " <br/> <br/> " + message

	return structures.Message{
		Message: cleanString,
	}
}
