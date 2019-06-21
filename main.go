package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc        string             `json:"desc,omitempty" bson:"desc,omitempty"`
	TimeCreated time.Time          `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Deadline    time.Time          `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Estimate    int                `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TimeSpent   int                `json:"timespent,omitempty" bson:"timespent,omitempty"`
}

var client *mongo.Client

// Creates a new Todo object in the database
func CreateTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	var todo Todo
	json.NewDecoder(req.Body).Decode(&todo)
	todo.TimeCreated = time.Now()
	todo.TimeSpent = 0

	collection := client.Database("gotodo").Collection("todos")

	result, _ := collection.InsertOne(context.TODO(), todo)

	json.NewEncoder(res).Encode(result)
}

func main() {
	fmt.Println("Starting...")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, _ = mongo.Connect(context.TODO(), clientOptions)

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	router := mux.NewRouter()
	router.HandleFunc("/", CreateTodoEndPoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))

}
