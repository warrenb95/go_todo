package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/warrenb95/go_todo/config"

	"github.com/warrenb95/go_todo/routes"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting...")

	fmt.Println("Connected to go todo on " + config.GotodoURL)

	router := mux.NewRouter()
	router.HandleFunc("/", routes.IndexEndPoint)
	router.HandleFunc("/newtodo", routes.CreateNewTodoEndPoint)
	router.HandleFunc("/{id}/delete", routes.DeleteTodoEndPoint)
	router.HandleFunc("/{id}", routes.TodoDetailEndPoint)
	router.HandleFunc("/{id}/update", routes.UpdateTodoEndPoint)
	router.HandleFunc("/{id}/timespent", routes.TimeSpentEndPoint)
	log.Fatal(http.ListenAndServe(":80", router))
}
