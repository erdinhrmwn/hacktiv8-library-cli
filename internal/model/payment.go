package model

import "time"

type Payment struct {
	ID        int
	InvoiceID int
	Amount    float64
	Date      time.Time
	Method    string

	Invoice *Invoice
}
