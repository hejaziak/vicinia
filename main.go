package main

import (
	"log"
	"net/http"
	"os"
	global "vicinia/globals"
	routes "vicinia/routes"

	cors "github.com/heppu/simple-cors"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := routes.NewRouter()
	global.InitSessions()
	if err:= global.InitMapClient("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo"; err != nil{
		pretty.Printf("fatal error: %s \n", err)		
	}

	log.Fatal(http.ListenAndServe(":"+port, cors.CORS(router)))
}
