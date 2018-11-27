package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Description string
}
type TodoCreation struct {
	Description string `json:"description,omitempty"`
}

var todos = []Todo{}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/todo", getAll).Methods("GET")
	router.HandleFunc("/todo", create).Methods("POST")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func create(w http.ResponseWriter, r *http.Request) {
	var data TodoCreation
	_ = json.NewDecoder(r.Body).Decode(&data)
	todos = append(todos, Todo{Description: data.Description})
	json.NewEncoder(w).Encode(data)

}
