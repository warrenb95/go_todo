package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const mongodbURL string = "http://localhost:3000/todo"

type Todo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc        string             `json:"desc,omitempty" bson:"desc,omitempty"`
	TimeCreated time.Time          `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Deadline    time.Time          `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Estimate    int64              `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TimeSpent   int64              `json:"timespent,omitempty" bson:"timespent,omitempty"`
}

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	allTodosRaw, err := http.Get(mongodbURL)
	if err != nil {
		// Handke error
	}

	fmt.Println(allTodosRaw)
}

func main() {
	fmt.Println("Starting...")

	fmt.Println("Connected to go todo!")

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	// router.HandleFunc("/", GetAllTodosEndPoint).Methods("GET")
	// router.HandleFunc("/{id}", GetTodoEndpoint).Methods("GET")
	// router.HandleFunc("/{id}", DeleteTodoEndPoint).Methods("DELETE")
	// router.HandleFunc("/{id}", UpdateTodoEndPoint).Methods("PUT")
	// router.HandleFunc("/{id}/timespent", TimeSpentEndPoint).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))

}
