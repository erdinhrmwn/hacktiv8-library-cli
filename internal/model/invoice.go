package model

import "time"

type Invoice struct {
	ID        int
	UserID    int
	LoanID    int
	Amount    float64
	IssueDate time.Time
	Status    string // "unpaid" | "paid"

	User *User
	Loan *Loan
}
