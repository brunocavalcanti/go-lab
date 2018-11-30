package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	db, err := gorm.Open("mysql", os.Getenv("MYSQL_URL"))
	defer db.Close()

	if err != nil {
		panic("failed to connect database")
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
		onError(w, "Invalid todo!", http.StatusUnauthorized)
		return
	}

	err = db.Create(&todo).Error

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

func createDataBase(schemas ...interface{}) {
	db.AutoMigrate(schemas)
	// for _, Item := range schemas {
	// 	db.HasTable(&Item)
	// 	db.CreateTable(&Item)

	// }
}
