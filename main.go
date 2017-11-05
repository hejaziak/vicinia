package main

import (
	"log"
	"net/http"
	"os"
	global "vicinia/globals"
	"vicinia/routes"

	cors "github.com/heppu/simple-cors"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		pretty.Println("Error loading .env file")
	}

	global.InitSessions()

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
