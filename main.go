package main

import (
	"log"
	"net/http"
	"os"
	datastructures "vicinia/datastructures"
	global "vicinia/globals"
	routes "vicinia/routes"

	cors "github.com/heppu/simple-cors"
)

import _ "github.com/joho/godotenv/autoload"

func main() {
	datastructures.InitSessions()
	datastructures.InitLocations()

	mapsKey := os.Getenv("GoogleMapsAPI")
	global.InitMapClient(mapsKey)

	aiKey := os.Getenv("AI")
	global.InitAiClient(aiKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := routes.NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, cors.CORS(router)))
}
