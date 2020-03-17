package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/go_todo/config"
	"github.com/warrenb95/go_todo/models"
)

// TodoDetailEndPoint to show the details of a Todo
func TodoDetailEndPoint(res http.ResponseWriter, req *http.Request) {
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

	err = config.Templates.ExecuteTemplate(res, "detail", todo)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
