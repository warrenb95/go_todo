package config

import (
	"html/template"
	"net/http"
)

// Templates obj for app.
var Templates = template.Must(template.ParseGlob("templates/*.html"))

// Url Strings for the app
const GotodoURL string = "http://localhost"
const MongodbURL string = "http://mongo_todo:3000/todo"

// Client http.client pointer
var Client = &http.Client{}
