package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"github.com/marcossegovia/apiai-go"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	welcomeMessage := structures.WelcomeStruct{
		Message: "Welcome ,where do you want to go ?",
		UUID:    uuid.NewV1(),
	}

	if err := json.NewEncoder(w).Encode(welcomeMessage); err != nil {
		log.Fatalf("fatal error: %s", err)
		returnError(w, "")
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

<<<<<<< HEAD
	inUUID, err := extractUUID(r)
	if err != nil {
		returnError(w, "sorry, UUID not set, please access \"/welcome\" to receive an UUID")
		return
	}

	var requestBody structures.Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
=======
	var requestBody structures.Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		panic(err)
	}

	uuid := extractUUID(r)
	index, err := strconv.Atoi(requestBody.Message)
	if err != nil {
		getList(w, r, uuid, requestBody.Message)
	} else {
		getDetails(w, r, uuid, index-1)
	}
}

func getList(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, message string) {
	c := global.GetMapClient()

	client, err := apiai.NewClient(
        &apiai.ClientConfig{
            Token:      "71027bbaf70a4a53847bedce6b83c94f",
            QueryLang:  "en",    //Default en
            SpeechLang: "en-US", //Default en-US
        },
    )
    if err != nil {
        fmt.Printf("%v", err)
    }
    //Set the query string and your current user identifier.
    qr, err := client.Query(apiai.Query{Query: []string{message}, SessionId: uuid.String()})
    if err != nil {
        fmt.Printf("%v", err)
    }
    action := qr.Result.Action

    if(strings.Compare(action,"search") != 0 ){
		if err := json.NewEncoder(w).Encode(qr.Result.Fulfillment.Speech); err != nil {
			panic(err)
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
			log.Fatalf("fatal error: %s", err)
		}

		output := SimplifyList(res.Results)

		jsonMessage, _ := json.Marshal(output)

		respondMessage := extractMessage(string(jsonMessage), "To get detailed information about a specific place, please type its ID")
		if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
			panic(err)
		}
		updateSession(extractUUID(r), res.Results)
    }
}

func extractUUID(r *http.Request) uuid.UUID {
	uuid, err := uuid.FromString(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Printf("Something gone wrong: %s", err)
	}

	return uuid
}

func updateSession(UUID uuid.UUID, input []maps.PlacesSearchResult) {
	placeIDs := make([]string, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			break
		}
		placeIDs[i] = input[i].PlaceID
	}

	global.UpdateEntry(UUID, placeIDs)
}

func getDetails(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, index int) {
	c := global.GetMapClient()

	placeID := global.GetPlace(uuid, index)

	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	res, err := c.PlaceDetails(context.Background(), req)
	if err != nil {
>>>>>>> e39bee6e1d11921f9a6fc6b1ba8ff5594067c53c
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

		body :=	"Available Routes:" +
			"  GET  /welcome - handleWelcome" +
			"  POST /chat    - handleChat" +
			"  GET  /        - handle        (current)" +" || "+ " We are happy to serve you !!"
	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}

}
