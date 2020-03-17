package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/go_todo/config"
)

// DeleteTodoEndPoint to delete a Todo from the database
func DeleteTodoEndPoint(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	// Set the correct URL for the Todo ID
	mongodbDeleteURL := config.MongodbURL + "/" + id

	// Create request
	deleteReq, err := http.NewRequest("DELETE", mongodbDeleteURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := config.Client.Do(deleteReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	http.Redirect(res, req, config.GotodoURL, http.StatusTemporaryRedirect)
}
