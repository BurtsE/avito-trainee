package models

import "time"

type User struct {
	Id        string
	UserName  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
