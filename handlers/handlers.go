package handlers

import (
	"encoding/json"
	"net/http"

	datastructures "vicinia/datastructures"
	helpers "vicinia/helperFunctions"
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

//IndexHandler : returns to the client the available routes
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := "Available Routes:<br/>" +
		"GET  /welcome - receive uuid <br/>" +
		"POST /chat    - converse with chatbot <br/>" +
		"GET  /        - welcome page        (current) <br/>" +
		"<br/>We are happy to serve you !!"

	helpers.ReturnMessage(w, response)
}

//WelcomeHandler : creates a new uuid and returns a welcome message
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newUUID := uuid.NewV1()
	response := structures.UUIDMessage{
		Message: "Welcome, where do you want to go today ?",
		UUID:    newUUID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		helpers.ReturnMessage(w, "")
		return
	}

	if err := datastructures.CreateEntry(newUUID); err != nil {
		pretty.Printf("fatal error: %s \n", err)
	}
}

//ChatHandler : returns the required chat message
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extract UUID
	inUUID, err := helpers.ExtractUUID(r)
	if err != nil {
		helpers.ReturnUnauthorized(w, "sorry, UUID not set, please access /welcome to receive an UUID")
		return
	}

	//ensure UUID is in database
	if _, err := datastructures.GetEntry(inUUID); err != nil {
		helpers.ReturnUnauthorized(w, "sorry, UUID is not correct, please access /welcome to receive a new UUID")
		return
	}

	//decode message json to object
	var requestBody structures.LatLongMessage
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		helpers.ReturnMessage(w, "")
		return
	}

	helpers.ListHandler(w, r, inUUID, requestBody)

}

//PlaceDetailsHandler: returns further details of a specific place
func PlaceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extract UUID
	inUUID, err := helpers.ExtractUUID(r)
	if err != nil {
		helpers.ReturnUnauthorized(w, "sorry, UUID not set, please access /welcome to receive an UUID")
		return
	}

	//ensure UUID is in database
	if _, err := datastructures.GetEntry(inUUID); err != nil {
		helpers.ReturnUnauthorized(w, "sorry, UUID is not correct, please access /welcome to receive a new UUID")
		return
	}

	placeID, ok := r.URL.Query()["place_id"]
    if !ok || len(placeID) < 1 {
		pretty.Printf("Url Param 'placeID' is missing")
		helpers.ReturnMessage(w, "")
		return
    }

	latitude, ok := r.URL.Query()["latitude"]
    if !ok || len(latitude) < 1 {
		pretty.Printf("Url Param 'latitude' is missing")
		helpers.ReturnMessage(w, "")
		return
    }

	longitude, ok := r.URL.Query()["longitude"]
    if !ok || len(longitude) < 1 {
		pretty.Printf("Url Param 'longitude' is missing")
		helpers.ReturnMessage(w, "")
		return
    }
	
	helpers.DetailsHandler(w, r, placeID[0], latitude[0], longitude[0])
}

// MiscHandler : handler for all unchaught routes
func MiscHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	helpers.ReturnMessage(w, "please perform a GET on access /welcome to get a new UUID or a POST on /chat to converse with the chatbot")
}
