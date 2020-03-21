package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/warrenb95/go_todo/config"
	"github.com/warrenb95/go_todo/models"
)

// TimeSpentEndPoint add timepent to a givent Todo
func TimeSpentEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	todoDetailURL := config.MongodbURL + "/" + id

	// PUT request URL
	timeSpentURL := todoDetailURL + "/timespent"

	// Create request
	getReq, err := http.NewRequest("GET", todoDetailURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := config.Client.Do(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todo models.Todo
	json.Unmarshal(resbodybytes, &todo)

	if req.Method == "GET" {
		err := config.Templates.ExecuteTemplate(res, "timespent", todo)
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

		updatedTimespent := append(todo.TimeSpent, models.Timespent{
			Duration: todoTimespent,
			Date:     time.Now(),
			Desc:     req.Form["description"][0],
		})

		var updatedtodo = models.Todo{
			TotalTimeSpent: todo.TotalTimeSpent + todoTimespent,
			TimeSpent:      updatedTimespent,
		}

		todojson, err := json.Marshal(updatedtodo)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Create request
		putReq, err := http.NewRequest("PUT", timeSpentURL, bytes.NewBuffer(todojson))
		if err != nil {
			fmt.Println(err)
			return
		}

		putReq.Header.Set("Content-Type", "application/json")

		// Fetch Request
		updateresp, err := config.Client.Do(putReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer updateresp.Body.Close()

		http.Redirect(res, req, config.GotodoURL+"/"+id, http.StatusTemporaryRedirect)
	}
}
