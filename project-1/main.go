package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type Todo struct {
	ID          uuid.UUID
	Description string
}
type TodoCreation struct {
	Description string `json:"description,omitempty"`
}

var todos = []Todo{}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/todo", getAll).Methods("GET")
	router.HandleFunc("/todo/{id}", get).Methods("GET")
	router.HandleFunc("/todo", create).Methods("POST")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}
func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidTodo := uuid.FromStringOrNil(vars["id"])

	for _, todo := range todos {
		if todo.ID == uuidTodo {
			json.NewEncoder(w).Encode(&todo)
			break
		}
	}
	json.NewEncoder(w).Encode(Todo{})

}
func create(w http.ResponseWriter, r *http.Request) {
	var data TodoCreation
	_ = json.NewDecoder(r.Body).Decode(&data)
	id, _ := uuid.NewV4()
	newTodo := Todo{Description: data.Description, ID: id}
	todos = append(todos, newTodo)
	json.NewEncoder(w).Encode(data)

}
