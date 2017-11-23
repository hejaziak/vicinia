package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
		Message: "Welcome, please enter yor location in the following format <br/> location:latitude,longitutde",
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
	var requestBody structures.Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		helpers.ReturnMessage(w, "")
		return
	}

	splittedMessage := strings.Split(requestBody.Message, ":")

	//forward message to appropriate handler
	switch splittedMessage[0] {
	case "location":
		if len(splittedMessage) == 1 { //handles this case--> "message":"location"
			helpers.ReturnMessage(w, "You have entered incorrect format for your location. please enter location:<latitude>,<longitude>")
			return
		}
		helpers.LocationHandler(w, r, inUUID, splittedMessage[1])
		break

	case "details":
		if len(splittedMessage) == 1 { //handles this case--> "message":"details"
			helpers.ReturnMessage(w, "You have entered incorrect format for your place details query. please enter details:<place index>")
			return
		}

		index, err := strconv.Atoi(splittedMessage[1])
		if err != nil {
			helpers.ReturnMessage(w, "You have entered incorrect format for your place details query. please enter details:<place index>")
			return
		}

		helpers.DetailsHandler(w, r, inUUID, index-1)
		break

	default:
		helpers.ListHandler(w, r, inUUID, requestBody.Message)
	}
}

// MiscHandler : handler for all unchaught routes
func MiscHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	helpers.ReturnMessage(w, "please perform a GET on access /welcome to get a new UUID or a POST on /chat to converse with the chatbot")
}
