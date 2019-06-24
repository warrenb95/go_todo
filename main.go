package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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
	Estimate    int64              `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TimeSpent   int64              `json:"timespent,omitempty" bson:"timespent,omitempty"`
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

func GetAllTodosEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	var todos []Todo
	collection := client.Database("gotodo").Collection("todos")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(todos)
}

func GetTodoEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	params := mux.Vars(req)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("gotodo").Collection("todos")

	var todo Todo
	err := collection.FindOne(context.TODO(), Todo{ID: id}).Decode(&todo)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(res).Encode(todo)
}

func DeleteTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("gotodo").Collection("todos")

	var todo Todo
	err := collection.FindOneAndDelete(context.TODO(), Todo{ID: id}).Decode(&todo)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(res).Encode(todo)
}

func UpdateTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("gotodo").Collection("todos")

	var todo Todo
	json.NewDecoder(req.Body).Decode(&todo)

	result, err := collection.UpdateOne(context.TODO(), Todo{ID: id}, bson.M{"$set": todo})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(res).Encode(result)
}

func TimeSpentEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("gotodo").Collection("todos")

	var todo Todo
	err := collection.FindOne(context.TODO(), Todo{ID: id}).Decode(&todo)

	var updatedTodo Todo
	json.NewDecoder(req.Body).Decode(&updatedTodo)

	todo.TimeSpent += updatedTodo.TimeSpent

	result, err := collection.UpdateOne(context.TODO(), Todo{ID: id}, bson.M{"$set": todo})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

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
	router.HandleFunc("/todo", CreateTodoEndPoint).Methods("POST")
	router.HandleFunc("/todo", GetAllTodosEndPoint).Methods("GET")
	router.HandleFunc("/todo/{id}", GetTodoEndpoint).Methods("GET")
	router.HandleFunc("/todo/{id}", DeleteTodoEndPoint).Methods("DELETE")
	router.HandleFunc("/todo/{id}", UpdateTodoEndPoint).Methods("PUT")
	router.HandleFunc("/todo/{id}/timespent", TimeSpentEndPoint).Methods("PUT")
	log.Fatal(http.ListenAndServe(":12345", router))

}
