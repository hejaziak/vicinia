package main

import (
	"log"
	"net/http"
	"os"
	global "vicinia/globals"
	routes "vicinia/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := routes.NewRouter()
	global.InitSessions()
	global.InitMapClient("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
