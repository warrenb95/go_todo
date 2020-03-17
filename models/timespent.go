package models

import (
	"fmt"
	"time"
)

// FormatAsDate to foramat the Date obj of Timespent
func (t Timespent) FormatAsDate() string {
	d := t.Date
	year, month, day := d.Date()
	return fmt.Sprintf("%d-%d-%d", day, month, year)
}

// Timespent struct for each todo obj
// Duration: Is int64 of the minutes spent
// Data: The date of the update
// Desc: The description of the update
type Timespent struct {
	Duration int64     `json:"timespent,omitempty" bson:"timespent,omitempty"`
	Date     time.Time `json:"timecreated,omitempty" bson:"timecreated,omitempty"`
	Desc     string    `json:"desc,omitempty" bson:"desc,omitempty"`
}
