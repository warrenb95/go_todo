package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const gotodoURL string = "http://localhost/"
const mongodbURL string = "http://localhost:3000/todo"

var templates = template.Must(template.ParseGlob("templates/*"))

type Todo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc        string             `json:"desc,omitempty" bson:"desc,omitempty"`
	TimeCreated time.Time          `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Deadline    time.Time          `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Estimate    int64              `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TimeSpent   int64              `json:"timespent,omitempty" bson:"timespent,omitempty"`
}

type IndexPageData struct {
	Todos []Todo
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

	data := IndexPageData{Todos: todos}

	err = templates.ExecuteTemplate(res, "index", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateNewTodoHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := templates.ExecuteTemplate(res, "newtodo", nil)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		req.ParseForm()

		estimate, err := strconv.ParseInt(req.Form["estimate"][0], 10, 64)
		if err == nil {
			fmt.Printf("%d of type %T", estimate, estimate)
		}

		var todo Todo = Todo{
			Title:     req.Form["title"][0],
			Desc:      req.Form["description"][0],
			Estimate:  estimate,
			TimeSpent: 0,
		}

		todojson, err := json.Marshal(todo)
		if err != nil {
			fmt.Println(err)
			return
		}

		createtodoresp, err := http.Post(mongodbURL, "application/json", bytes.NewBuffer(todojson))
		if err != nil {
			log.Fatalln(err)
		}

		createtodoresp.Body.Close()

		http.Redirect(res, req, gotodoURL, http.StatusSeeOther)
	}
}

func main() {
	fmt.Println("Starting...")

	fmt.Println("Connected to go todo!")

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/newtodo", CreateNewTodoHandler)
	// router.HandleFunc("/{id}", TodoDetail)
	// router.HandleFunc("/{id}", DeleteTodoEndPoint)
	// router.HandleFunc("/{id}", UpdateTodoEndPoint)
	// router.HandleFunc("/{id}/timespent", TimeSpentEndPoint)
	log.Fatal(http.ListenAndServe(":80", router))

}
