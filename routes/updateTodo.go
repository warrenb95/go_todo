package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/warrenb95/go_todo/config"
	"github.com/warrenb95/go_todo/models"
)

// UpdateTodoEndPoint to update the Todo
func UpdateTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	todoDetailURL := config.MongodbURL + "/" + id

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
		err := config.Templates.ExecuteTemplate(res, "update", todo)
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

		var updatedtodo = models.Todo{
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
		updateresp, err := config.Client.Do(putReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer updateresp.Body.Close()

		http.Redirect(res, req, config.GotodoURL+"/"+id, http.StatusTemporaryRedirect)
	}
}
