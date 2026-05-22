package model

import "time"

type Loan struct {
	ID         int
	VisitorID  int
	StaffID    int
	BookID     int
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string // "active" | "returned"

	Visitor *User
	Staff   *User
	Book    *Book
}
