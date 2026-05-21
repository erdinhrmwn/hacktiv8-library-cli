package model

import "time"

type Loan struct {
	ID         int
	VisitorID  int
	StaffID    int
	BookID     int
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate *time.Time // pointer → nullable
	Status     string     // "active" | "returned"

	// Relasi — diisi manual lewat JOIN di repository
	Visitor *User
	Book    *Book
}
