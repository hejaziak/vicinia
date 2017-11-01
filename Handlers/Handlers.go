package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	s "Vicinia/Structures"

)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := s.Todos{
		s.Todo{Name: "Write presentation"},
		s.Todo{Name: "Host meetup"},
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

	welcomeMessage := s.WelcomeStruct{
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

	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}
