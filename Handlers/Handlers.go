package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	structures "Vicinia/Structures"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := structures.Todos{
		structures.Todo{Name: "Write presentation"},
		structures.Todo{Name: "Host meetup"},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {

	welcomeMessage := structures.WelcomeStruct{
		Message: "Welcome ,where do you want to go ?",
		Uuid:    uuid.NewV1(),
	}

	if err := json.NewEncoder(w).Encode(welcomeMessage); err != nil {
		panic(err)
	}

}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: 29.985352, Lng: 31.279194},
		RankBy:   "distance",
		Keyword:  "resturants",
	}

	res, err := c.NearbySearch(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	output := SimplifyList(res.Results)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		panic(err)
	}
}

func SimplifyList(input []maps.PlacesSearchResult) []structures.PlaceListEntity {
	output := make([]structures.PlaceListEntity, 5)

	for i := 0; i < 5; i++ {
		if i >= len(input) {
			break
		}

		c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"))
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		req := &maps.DistanceMatrixRequest{
			Origins:      []string{"29.985352,31.279194"},
			Destinations: []string{"place_id:" + input[i].PlaceID},
			Units:        "metric",
		}

		res, err := c.DistanceMatrix(context.Background(), req)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		pretty.Println(res)

		output[i] = structures.PlaceListEntity{
			Name:     input[i].Name,
			Distance: res.Rows[0].Elements[0].Distance.HumanReadable,
			Rating:   input[i].Rating,
		}
	}

	return output
}
