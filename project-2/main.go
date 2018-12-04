package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Todo struct {
	gorm.Model
	Description string
}
type TodoCreation struct {
	Description string `json:"description,required,omitempty"`
}
type ErrorHandlerApi struct {
	Message string `json:"message,omitempty"`
}

var db *gorm.DB
var todos = []Todo{}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "./todo.db")
	db.AutoMigrate(&Todo{})
	defer db.Close()

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/todo", getAll).Methods("GET")
	router.HandleFunc("/todo/{id}", get).Methods("GET")
	router.HandleFunc("/todo", create).Methods("POST")
	router.HandleFunc("/todo/{id}", remove).Methods("DELETE")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	onSucces(w, todos, http.StatusOK)
}
func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)

}

func create(w http.ResponseWriter, r *http.Request) {
	var todo TodoCreation
	var err error
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		onError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (TodoCreation{}) == todo {
		onError(w, err.Error(), http.StatusBadRequest)
		return
	}

	retInsert := db.Create(&todo)
	if retInsert.Error != nil {
		onError(w, "TESTE", http.StatusBadRequest)
		return
	}
	onSucces(w, retInsert.Value, http.StatusCreated)
	return

}

func remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
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

func createDataBase() {
	// db.AutoMigrate(&Todo{})
	// for _, Item := range schemas {
	// 	db.HasTable(&Item)
	// 	db.CreateTable(&Item)

	// }
}
