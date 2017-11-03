package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	globals "vicinia/globals"
	structures "vicinia/structures"

	"github.com/satori/go.uuid"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := "Available Routes:" +
		"  GET  /welcome - handleWelcome" +
		"  POST /chat    - handleChat" +
		"  GET  /        - handle        (current)" + " || " + " We are happy to serve you !!"

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatalf("fatal error: %s", err)
		returnError(w, "")
	}
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newUUID := uuid.NewV1()
	welcomeMessage := structures.WelcomeStruct{
		Message: "Welcome ,where do you want to go ?",
		UUID:    newUUID,
	}

	if err := json.NewEncoder(w).Encode(welcomeMessage); err != nil {
		log.Fatalf("fatal error: %s", err)
		returnError(w, "")
	}

	if err := globals.CreateEntry(newUUID); err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	inUUID, err := extractUUID(r)
	if err != nil {
		returnError(w, "sorry, UUID not set, please access \"/welcome\" to receive an UUID")
		return
	}

	var requestBody structures.Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Fatalf("fatal error: %s", err)
		returnError(w, "")
	}

	index, err := strconv.Atoi(requestBody.Message)
	if err != nil {
		getList(w, r, inUUID, requestBody.Message)
	} else {
		getDetails(w, r, inUUID, index-1)
	}
}
