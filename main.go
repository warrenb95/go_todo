package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	resbody, err := http.Get(mongodbURL)
	if err != nil {
		// Handke error
	}
	defer resbody.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resbody.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todos []Todo
	json.Unmarshal(resbodybytes, &todos)

}

func main() {
	fmt.Println("Starting...")

	fmt.Println("Connected to go todo!")

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	// router.HandleFunc("/{id}", TodoDetail)
	// router.HandleFunc("/{id}", DeleteTodoEndPoint)
	// router.HandleFunc("/{id}", UpdateTodoEndPoint)
	// router.HandleFunc("/{id}/timespent", TimeSpentEndPoint)
	log.Fatal(http.ListenAndServe(":80", router))

}
