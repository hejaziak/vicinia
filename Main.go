package main

import (
    "log"
    "net/http"
    r "Vicinia/Routes"
)

func main() {

    router := r.NewRouter()

    log.Fatal(http.ListenAndServe(":8080", router))
}