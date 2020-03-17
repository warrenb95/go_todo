package config

import (
	"html/template"
	"net/http"
)

// Templates obj for app.
var Templates = template.Must(template.ParseGlob("templates/*.html"))

// Url Strings for the app
const GotodoURL string = "http://localhost/"
const MongodbURL string = "http://localhost:3000/todo"

// The http.client pointer
// var client *http.Client = &http.Client{}
var Client = &http.Client{}
