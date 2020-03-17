package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/warrenb95/go_todo/config"

	"github.com/warrenb95/go_todo/models"
)

// CreateNewTodoEndPoint to create a new Todo
func CreateNewTodoEndPoint(res http.ResponseWriter, req *http.Request) {

	// Handle the 'GET' and 'POST' methods
	if req.Method == "GET" {
		// Execute the 'newtodo' template
		err := config.Templates.ExecuteTemplate(res, "newtodo", nil)
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
		var todo = models.Todo{
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
		createtodoreq, err := http.NewRequest("POST", config.MongodbURL, bytes.NewBuffer(todojson))
		if err != nil {
			log.Fatalln(err)
		}
		createtodoreq.Header.Set("Content-Type", "application/json")

		resp, err := config.Client.Do(createtodoreq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		http.Redirect(res, req, config.GotodoURL, http.StatusTemporaryRedirect)
	}
}
