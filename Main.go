package main

import (
	global "Vicinia/Globals"
	routes "Vicinia/Routes"
	"log"
	"net/http"
)

func main() {
	router := routes.NewRouter()
	global.InitSessions()
	global.InitMapClient("AIzaSyBZwHSODUVFhzMcAEabT-BOw2_SkOrYEWo")

	log.Fatal(http.ListenAndServe(":8080", router))
}
