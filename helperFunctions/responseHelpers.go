package helperFunctions

import (
	"encoding/json"
	"net/http"
	"strings"

	structures "vicinia/structures"

	"github.com/kr/pretty"
	uuid "github.com/satori/go.uuid"
)

//ListHandler : returns the first 5 nearby places obtained from Google Maps API and updates the session map
//with the current places returned to the user
func ListHandler(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, message string) {
	aiResponse, err := GetIntent(uuid, message)
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

	keyword := aiResponse.Params["keyword"]
	result, err := GetList(uuid, keyword)
	if err != nil {
		if len(result) == 0 {
			ReturnMessage(w, "Couldn't find any nearby places with keyword specified")
		} else {
			ReturnMessage(w, "")
		}

		return
	}

	output, err := SimplifyList(uuid, result)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

	respondMessage := formatList(output, "To get detailed information about a specific place, please type its ID")
	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}

	UpdateSession(uuid, result)
}

func DetailsHandler(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, index int) {
	result, err := GetDetails(uuid, index)
	if err != nil {
		ReturnMessage(w, "")
		return
	}

	output, err := SimplifyDetails(uuid, result)

	respondMessage := formatDetails(output, "Any other place you want to search for ?")
	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "")
		return
	}
}

func LocationHandler(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, location string) {
	if err := CheckLocationFormat(location); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		ReturnMessage(w, "You have entered incorrect format for your location. please enter location:<latitude>,<longitude>")
		return
	}

	reply := SetLocation(uuid, location)
	if reply != "" {
		ReturnMessage(w, reply)
		return
	}

	ReturnMessage(w, "You're location is now set, where do you want to go today ?")
}

//ReturnUnauthorized : returns unauthorized error message
func ReturnUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)

	ReturnMessage(w, message)
}

//ReturnMessage : returns a message
func ReturnMessage(w http.ResponseWriter, message string) {
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
