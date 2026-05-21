package model

import "time"

type Invoice struct {
	ID        int
	UserID    int // visitor yang kena denda
	LoanID    int
	Amount    float64
	IssueDate time.Time
	Status    string // "unpaid" | "paid"
}
