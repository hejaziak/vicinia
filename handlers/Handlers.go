package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	global "vicinia/globals"
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
		panic(err)
	}

	global.CreateEntry(welcomeMessage.UUID)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		log.Fatalf("fatal error: %s", err)
	}

	output := SimplifyDetails(res)
	pretty.Println(res)

	jsonMessage, _ := json.Marshal(output)

	respondMessage := extractMessage(string(jsonMessage), "")
	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
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
			ID:       i + 1,
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
	c := global.GetMapClient()

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

func extractMessage(json string, message string) structures.Message {
	s2 := strings.Replace(json, "{", "", -1)
	s3 := strings.Replace(s2, "}", "", -1)
	s4 := strings.Replace(s3, "[", "", -1)
	s5 := strings.Replace(s4, "]", "", -1)
	s6 := strings.Replace(s5, "\"", "", -1)
	cleanString := strings.Replace(s6, ",", " <br/> ", -1)
	cleanString = cleanString + " <br/> " + message

	pretty.Println(message)
	return structures.Message{
		Message: cleanString,
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
