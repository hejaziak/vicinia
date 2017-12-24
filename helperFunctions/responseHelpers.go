package helperFunctions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	structures "vicinia/structures"

	"github.com/kr/pretty"
	uuid "github.com/satori/go.uuid"
)

//ListHandler : returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func ListHandler(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, requestBody structures.LatLongMessage) {
	aiResponse, err := GetIntent(uuid, requestBody.Message)
	if err != nil {
		ReturnMessage(w, "")
		return
	}

	//if action != search, that means the AI handled the response
	action := aiResponse.Action
	if strings.Compare(action, "search") != 0 {
		ReturnMessage(w, aiResponse.Fulfillment.Speech)
		return
	}

	latitude, err := strconv.ParseFloat(requestBody.Latitude, 64)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

	longitude, err := strconv.ParseFloat(requestBody.Longitude, 64)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

	keyword := aiResponse.Params["keyword"].(string)
	result, err := GetList(latitude, longitude, keyword)
	if err != nil {
		if len(result) == 0 {
			ReturnMessage(w, "Couldn't find any nearby places with keyword specified")
		} else {
			ReturnMessage(w, "")
		}

		return
	}

	output, err := SimplifyList(requestBody.Latitude, requestBody.Longitude, result)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

	response := structures.PlaceListMessage{
		Message: output,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

}

//DetailsHandler : returns details about a place, given the placeID
func DetailsHandler(w http.ResponseWriter, r *http.Request, placeID string, latitude string, longitude string) {

	result, err := GetDetails(placeID)
	if err != nil {
		ReturnMessage(w, "")
		return
	}

	output, err := SimplifyDetails(latitude, longitude, result)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

}

//ReturnMessage : returns a message
func ReturnMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")

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

//ReturnError : returns 400 error message
func ReturnError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)

	ReturnMessage(w, message)
}

//ReturnUnauthorized : returns unauthorized error message
func ReturnUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)

	ReturnMessage(w, message)
}
