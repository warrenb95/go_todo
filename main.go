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

// Url Strings for the app
const gotodoURL string = "http://localhost/"
const mongodbURL string = "http://localhost:3000/todo"

// Templates obj for app.
var templates = template.Must(template.ParseGlob("templates/*.html"))

// Timespent struct for each todo obj
// Duration: Is int64 of the minutes spent
// Data: The date of the update
// Desc: The description of the update
type Timespent struct {
	Duration int64     `json:"timespent,omitempty" bson:"timespent,omitempty"`
	Date     time.Time `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Desc     string    `json:"desc,omitempty" bson:"desc,omitempty"`
}

// FormatAsDate to foramat the Date obj of Timespent
func (t Timespent) FormatAsDate() string {
	d := t.Date
	year, month, day := d.Date()
	return fmt.Sprintf("%d-%d-%d", day, month, year)
}

// Todo struct to hold values for each Todo
// ID: unique ID for all Todo's
// Title: Title of Todo as string
// Desc: Description on the Todo as a string
// TimeCreated: The time the Todo was created time.Time
// Deadline: The deadline of the Todo as time.Time
// Estimate: Minutes as int64
// TotalTimeSpent: Minuets as int64
// TimeSpent: List of TimeSpent structs
type Todo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc           string             `json:"desc,omitempty" bson:"desc,omitempty"`
	TimeCreated    time.Time          `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Deadline       time.Time          `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Estimate       int64              `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TotalTimeSpent int64              `json:"totaltimespent,omitempty" bson:"totaltimespent,omitempty"`
	TimeSpent      []Timespent        `json:"timespent,omitempty" bson:"timespent,omitempty"`
}

// IndexPageData the list of Todos to display on the index page
type IndexPageData struct {
	Todos []Todo
}

// The http.client pointer
var client *http.Client = &http.Client{}

// IndexEndPoint for the home/index page
func IndexEndPoint(res http.ResponseWriter, req *http.Request) {
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

// CreateNewTodoEndPoint to create a new Todo
func CreateNewTodoEndPoint(res http.ResponseWriter, req *http.Request) {

	// Handle the 'GET' and 'REQUEST' methods
	if req.Method == "GET" {
		// Execute the 'newtodo' template
		err := templates.ExecuteTemplate(res, "newtodo", nil)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Parse the 'newtodo' form
		req.ParseForm()

		estimate, err := strconv.ParseInt(req.Form["estimate"][0], 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", estimate, estimate)
		}

		// Create a new Todo and set its values
		var todo = Todo{
			Title:          req.Form["title"][0],
			Desc:           req.Form["description"][0],
			Estimate:       estimate,
			TotalTimeSpent: 0,
		}

		// Marshel the Todo into json
		todojson, err := json.Marshal(todo)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Post the new Todo to the mondo database api
		createtodoreq, err := http.NewRequest("POST", mongodbURL, bytes.NewBuffer(todojson))
		if err != nil {
			log.Fatalln(err)
		}
		createtodoreq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(createtodoreq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		http.Redirect(res, req, gotodoURL, http.StatusTemporaryRedirect)
	}
}

// DeleteTodoEndPoint to delete a Todo from the database
func DeleteTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	mongodbDeleteURL := mongodbURL + "/" + id

	// Create request
	deleteReq, err := http.NewRequest("DELETE", mongodbDeleteURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(deleteReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	http.Redirect(res, req, gotodoURL, http.StatusTemporaryRedirect)
}

// TodoDetailEndPoint to show the details of a Todo
func TodoDetailEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	todoDetailURL := mongodbURL + "/" + id

	// Create request
	getReq, err := http.NewRequest("GET", todoDetailURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todo Todo
	json.Unmarshal(resbodybytes, &todo)

	err = templates.ExecuteTemplate(res, "detail", todo)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateTodoEndPoint to update the Todo
func UpdateTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	todoDetailURL := mongodbURL + "/" + id

	// Create request
	getReq, err := http.NewRequest("GET", todoDetailURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todo Todo
	json.Unmarshal(resbodybytes, &todo)

	if req.Method == "GET" {
		err := templates.ExecuteTemplate(res, "update", todo)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		req.ParseForm()

		estimate, err := strconv.ParseInt(req.Form["estimate"][0], 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", estimate, estimate)
		}

		var updatedtodo = Todo{
			Title:    req.Form["title"][0],
			Desc:     req.Form["description"][0],
			Estimate: estimate,
		}

		todojson, err := json.Marshal(updatedtodo)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Create request
		putReq, err := http.NewRequest("PUT", todoDetailURL, bytes.NewBuffer(todojson))
		if err != nil {
			fmt.Println(err)
			return
		}

		putReq.Header.Set("Content-Type", "application/json")

		// Fetch Request
		updateresp, err := client.Do(putReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer updateresp.Body.Close()

		http.Redirect(res, req, gotodoURL+"/"+id, http.StatusTemporaryRedirect)
	}
}

// TimeSpentEndPoint add timepent to a givent Todo
func TimeSpentEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	todoDetailURL := mongodbURL + "/" + id

	// Create request
	getReq, err := http.NewRequest("GET", todoDetailURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todo Todo
	json.Unmarshal(resbodybytes, &todo)

	if req.Method == "GET" {
		err := templates.ExecuteTemplate(res, "timespent", todo)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		req.ParseForm()

		todoTimespent, err := strconv.ParseInt(req.Form["timespent"][0], 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", todoTimespent, todoTimespent)
		}

		updatedTimespent := append(todo.TimeSpent, Timespent{Duration: todoTimespent,
			Date: time.Now(),
			Desc: req.Form["description"][0]})

		var updatedtodo = Todo{
			TotalTimeSpent: todo.TotalTimeSpent + todoTimespent,
			TimeSpent:      updatedTimespent,
		}

		todojson, err := json.Marshal(updatedtodo)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Create request
		putReq, err := http.NewRequest("PUT", todoDetailURL, bytes.NewBuffer(todojson))
		if err != nil {
			fmt.Println(err)
			return
		}

		putReq.Header.Set("Content-Type", "application/json")

		// Fetch Request
		updateresp, err := client.Do(putReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer updateresp.Body.Close()

		http.Redirect(res, req, gotodoURL+"/"+id, http.StatusTemporaryRedirect)
	}
}

// DisplayOverviewEndPoint nto sure what this does oops
func DisplayOverviewEndPoint(res http.ResponseWriter, req *http.Request) {
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

	err = templates.ExecuteTemplate(res, "overview", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Starting...")

	fmt.Println("Connected to go todo!")

	router := mux.NewRouter()
	router.HandleFunc("/", IndexEndPoint)
	router.HandleFunc("/newtodo", CreateNewTodoEndPoint)
	router.HandleFunc("/{id}/delete", DeleteTodoEndPoint)
	router.HandleFunc("/{id}", TodoDetailEndPoint)
	router.HandleFunc("/{id}/update", UpdateTodoEndPoint)
	router.HandleFunc("/{id}/timespent", TimeSpentEndPoint)
	log.Fatal(http.ListenAndServe(":80", router))
}
