package model

import "time"

type Author struct {
	ID          int
	Name        string
	BirthDate   time.Time
	Nationality string
	Bio         string
}
