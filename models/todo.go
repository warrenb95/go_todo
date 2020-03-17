package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo struct to hold values for each Todo
// ID: unique ID for all Todo's
// Title: Title of Todo as string
// Desc: Description on the Todo as a string
// TimeCreated: The time the Todo was created time.Time
// Deadline: The deadline of the Todo as time.Time
// Estimate: Minutes as int64
// TotalTimeSpent: Minuets as int64
// TimeSpent: List of TimeSpent structs
type Todo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc           string             `json:"desc,omitempty" bson:"desc,omitempty"`
	TimeCreated    time.Time          `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Deadline       time.Time          `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Estimate       int64              `json:"estimate,omitempty" bson:"estimate,omitempty"`
	TotalTimeSpent int64              `json:"totaltimespent,omitempty" bson:"totaltimespent,omitempty"`
	TimeSpent      []Timespent        `json:"timespent,omitempty" bson:"timespent,omitempty"`
}
