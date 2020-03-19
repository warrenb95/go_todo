package config

import (
	"html/template"
	"net/http"
)

// Templates obj for app.
var Templates = template.Must(template.ParseGlob("templates/*.html"))

// GotodoURL Url for the app
const GotodoURL string = "http://go_todo:80"

// MongodbURL for the mongo API
const MongodbURL string = "http://mongo_todo:3000/todo"

// const MongodbURL string = "http://172.28.0.3:3000/todo"

// Client http.client pointer
var Client = &http.Client{}
