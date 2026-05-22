package model

import "time"

type ActivityLog struct {
	ID          int
	Key         string
	Description string
	Date        time.Time
}
