package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"

    "github.com/satori/go.uuid"
    
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    todos := Todos{
        Todo{Name: "Write presentation"},
        Todo{Name: "Host meetup"},
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

    welcomeMessage := WelcomeStruct{
        Message:"Welcome ,where do you want to go ?",
        Uuid: uuid.NewV1(),
    }

    if err := json.NewEncoder(w).Encode(welcomeMessage); err != nil {
        panic(err)
    }
    

}