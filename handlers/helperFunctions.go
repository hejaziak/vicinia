package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	global "vicinia/globals"
	structures "vicinia/structures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"github.com/marcossegovia/apiai-go"
)

func getList(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, message string) {
	c, err := global.GetMapClient()
	if err != nil {
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
        pretty.Printf("%v", err)
    }

    action := qr.Result.Action

	if(strings.Compare(action,"search") != 0 ){

		res := structures.Message{
			Message: qr.Result.Fulfillment.Speech,
		} 
		if err := json.NewEncoder(w).Encode(res); err != nil {
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
			returnError(w, "")
		}

		if len(res.Results) <= 0 {
			returnError(w, "sorry I couldn't find any results matching the keyword: "+string(keyword.(string)))
			return
		}
		output , err:= SimplifyList(res.Results)

		jsonMessage, _ := json.Marshal(output)

		respondMessage := extractMessage(string(jsonMessage), "To get detailed information about a specific place, please type its ID")
		if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		inUUID, err := extractUUID(r)
		if err != nil {
			returnError(w, "")
			return
		}
		updateSession(inUUID, res.Results)
    }
}

func getDetails(w http.ResponseWriter, r *http.Request, uuid uuid.UUID, index int) {
	c, err := global.GetMapClient()
	if err != nil {
		returnError(w, "")
		return
	}

	placeID, err := global.GetPlace(uuid, index)
	if err != nil {
		returnError(w, "")
		return
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	res, err := c.PlaceDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	output, err := SimplifyDetails(res)
	if err != nil {
		returnError(w, "")
		return
	}
	pretty.Println(res)

	jsonMessage, _ := json.Marshal(output)
	respondMessage := extractMessage(string(jsonMessage), "Any other place you want to search for ?")
	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func extractUUID(r *http.Request) (uuid.UUID, error) {
	uuidNew, err := uuid.FromString(r.Header.Get("Authorization"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
		return uuid.Nil, err
	}

	return uuidNew, nil
}

func updateSession(UUID uuid.UUID, input []maps.PlacesSearchResult) error {
	placeIDs := make([]string, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			break
		}
		placeIDs[i] = input[i].PlaceID
	}

	err := global.UpdateEntry(UUID, placeIDs)
	if err != nil {
		return err
	}

	return nil
}

func returnError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Oops! something went wrong"
	}

	respondMessage := structures.Message{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(respondMessage); err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

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

		distance, err := getDistance("29.985352,31.279194", input[i].PlaceID)
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

func SimplifyDetails(input maps.PlaceDetailsResult) (structures.Place, error) {
	name := input.Name
	if name == "" {
		name = "not specified"
	}

	distance, err := getDistance("29.985352,31.279194", input.PlaceID)
	if err != nil {
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

func getDistance(cord string, destination string) (string, error) {
	c, err := global.GetMapClient()
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
