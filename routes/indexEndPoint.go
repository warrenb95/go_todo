package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/warrenb95/go_todo/config"
	"github.com/warrenb95/go_todo/models"
)

// IndexPageData the list of Todos to display on the index page
type IndexPageData struct {
	Todos []models.Todo
}

// IndexEndPoint for the home/index page
func IndexEndPoint(res http.ResponseWriter, req *http.Request) {
	resbody, err := config.Client.Get(config.MongodbURL)
	fmt.Println(resbody.Header)
	if err != nil {
		log.Fatal(err)
	}
	defer resbody.Body.Close()

	resbodybytes, err := ioutil.ReadAll(resbody.Body)
	if err != nil {
		log.Fatal(err)
	}

	var todos []models.Todo
	json.Unmarshal(resbodybytes, &todos)

	data := IndexPageData{Todos: todos}

	err = config.Templates.ExecuteTemplate(res, "index", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
