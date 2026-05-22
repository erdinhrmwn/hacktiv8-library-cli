package model

import "time"

type ActivityLog struct {
	ID          int
	Title       string
	Description string
	Date        time.Time
}
