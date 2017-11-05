package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	globals "vicinia/globals"
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

//IndexHandler: returns to the client the available routes
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := "Available Routes:" +
		"  GET  /welcome - handleWelcome" +
		"  POST /chat    - handleChat" +
		"  GET  /        - handle        (current)" + " || " + " We are happy to serve you !!"

	if err := json.NewEncoder(w).Encode(body); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}
}

//WelcomeHandler: creates a new uuid and returns a welcome message
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newUUID := uuid.NewV1()
	response := structures.WelcomeStruct{
		Message: "Welcome ,where do you want to go ?",
		UUID:    newUUID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	if err := globals.CreateEntry(newUUID); err != nil {
		pretty.Printf("fatal error: %s \n", err)
	}
}

//ChatHandler: returns the required chat message
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	inUUID, err := extractUUID(r)
	if err != nil {
		returnUnauthorized(w, "sorry, UUID not set, please access /welcome to receive an UUID")
		return
	}

	if _, err := globals.GetEntry(inUUID); err != nil {
		returnUnauthorized(w, "sorry, UUID is not correct, please access /welcome to receive an UUID")
		return
	}

	var requestBody structures.Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		returnError(w, "")
		return
	}

	index, err := strconv.Atoi(requestBody.Message)
	if err != nil {
		getList(w, r, inUUID, requestBody.Message)
	} else {
		getDetails(w, r, inUUID, index-1)
	}
}

func MiscHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	returnError(w, "please perform a GET on access /welcome to get a new UUID or a POST on /chat to converse with the chatbot")
}
