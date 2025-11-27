package models

import "time"

type Employee struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Password    string    `json:"-"`
	Designation string    `json:"designation"`
	CreatedAt   time.Time `json:"created_at"`
}
