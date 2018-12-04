package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"uuid,omitempty"`
	Description string    `json:"description,omitempty"`
}
type TodoCreation struct {
	Description string `json:"description,omitempty"`
}
type ErrorHandlerApi struct {
	Message string `json:"message,omitempty"`
}

var todos = []Todo{}

func main() {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/todo", getAll).Methods("GET")
	router.HandleFunc("/todo/{id}", get).Methods("GET")
	router.HandleFunc("/todo", create).Methods("POST")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	onSucces(w, todos, http.StatusOK)
}
func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuidTodo, err := uuid.FromString(vars["id"])

	if err != nil {
		onError(w, "Invalid Invalid!", http.StatusBadRequest)
		return
	}
	for _, item := range todos {
		if item.ID == uuidTodo {
			onSucces(w, item, http.StatusOK)
			return
		}
	}
	onError(w, "Todo not found!", http.StatusNotFound)
}

func create(w http.ResponseWriter, r *http.Request) {
	var data TodoCreation
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		onError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := uuid.NewV4()
	newTodo := Todo{Description: data.Description, ID: id}
	todos = append(todos, newTodo)

	onSucces(w, newTodo, http.StatusCreated)
}
func onError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	err := ErrorHandlerApi{message}
	json.NewEncoder(w).Encode(err)
}
func onSucces(w http.ResponseWriter, ret interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ret)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
