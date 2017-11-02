package main

import (
	datastructures "Vicinia/DataStructures"
	routes "Vicinia/Routes"
	"log"
	"net/http"
)

func main() {

	router := routes.NewRouter()
	datastructures.Init()
	log.Fatal(http.ListenAndServe(":8080", router))
}
