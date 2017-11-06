package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"strconv"
	
	globals "vicinia/globals"
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/marcossegovia/apiai-go"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

//setLocation: sets the location of client in locations map
func setLocation(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, location string){
	coordinates := strings.Split(location , ",")

	if err:=globals.CreateLocationEntry(uuid, coordinates[0], coordinates[1]); err!=nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "You have already entered your location successfully")
		return
	}

	respondMessage :=structures.Message{
		Message: "That's great.Now, where do you want to go today ?",
	} 
	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	} 
}

//getList: returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func getList(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, message string) {
	//getting the Google Maps client
	client, err := globals.GetAiClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	//Set the query string and your current user identifier.
	qr, err := client.Query(apiai.Query{
		Query:     []string{message},
		SessionId: uuid.String(),
	})

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	action := qr.Result.Action

	if strings.Compare(action, "search") != 0 {
		respondMessage := structures.Message{
			Message: qr.Result.Fulfillment.Speech,
		}
		if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

	} else {

		//getting the Google Maps client
		mapsClient, err := globals.GetMapClient()
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

		keyword := qr.Result.Params["keyword"]

		location,err := globals.GetLocationEntry(uuid) //location contains latitude and longitude
		if err!=nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "Please provide me first with your location")
			return			
		}

		latitude, err := strconv.ParseFloat(location[0], 64)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

		longitude, err := strconv.ParseFloat(location[1], 64)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}

		req := &maps.NearbySearchRequest{
			Location: &maps.LatLng{Lat: latitude, Lng: longitude},
			RankBy:   "distance",
			Keyword:  string(keyword.(string)),
		}

		res, err := mapsClient.NearbySearch(context.Background(), req)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "sorry, I couldn't find any relevant places")
			return
		}
		
		output, err := SimplifyList(uuid, res.Results)
		if err != nil {
			pretty.Printf("fatal error: %s \n", err)
			returnError(w, "")
			return
		}
		
		respondMessage := formatList(output, "To get detailed information about a specific place, please type its ID")
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

	output, err := SimplifyDetails(uuid, res)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	respondMessage := formatDetails(output, "Any other place you want to search for ?")

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
func SimplifyList(uuid uuid.UUID, input []maps.PlacesSearchResult) ([]structures.PlaceListEntity, error) {
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

		location, err := globals.GetLocationEntry(uuid) //location contains latitude and longitude
		if err!=nil {
			pretty.Printf("fatal error: %s \n", err)
			return []structures.PlaceListEntity{}, err			
		}
		
		distance, err := getDistance(location[0]+","+location[1], input[i].PlaceID) 
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

//SimplifyDetails : returns the list of parameters which will be returned as a response message to the user's specific query
func SimplifyDetails(uuid uuid.UUID, input maps.PlaceDetailsResult) (structures.Place, error) {
	name := input.Name
	if name == "" {
		name = "not specified"
	}

	location, err := globals.GetLocationEntry(uuid) //location contains latitude and longitude
	if err!=nil {
		pretty.Printf("fatal error: %s \n", err)
		return structures.Place{}, err			
	}

	distance, err := getDistance(location[0]+","+location[1], input.PlaceID) 
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

//formatList: returns a formatted message containing list of places
func formatList(placesList []structures.PlaceListEntity, message string) structures.Message{
	formattedMessage := ""
	for _, place := range placesList{
		formattedMessage+=
		"Name: " + place.Name + " <br/> " +
		"Distance: " + place.Distance + " <br/> " +
		"Rating: " + strconv.FormatFloat(float64(place.Rating), 'f', -1, 32) + " <br/> " +
		"ID: " + strconv.Itoa(place.ID) + " <br/> <br/> "
	}

	formattedMessage+=message

	return structures.Message{
		Message: formattedMessage,
	}
}

//formatDetails: returns a formatted message containing details of a specifice place
func formatDetails(placeDetails structures.Place, message string) structures.Message{
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

//checkLocationFormat: checks that the location is entered in the correct format
func checkLocationFormat(location string) error{
	
	if location=="" {
		return errors.New("empty location")
	}
	splittedLocation := strings.Split(location , ",")
	if len(splittedLocation) > 2{
		return errors.New("more than one ',' in location")
	}
	if _, err := strconv.ParseFloat(splittedLocation[0], 64); err!=nil{
		return errors.New("syntax error in latitude")
	}
	if _, err := strconv.ParseFloat(splittedLocation[1], 64); err!=nil{
		return errors.New("syntax error in longitude")
	}

	return nil
}
