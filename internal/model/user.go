package model

// User represents a row in the users table.
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string // "staff" | "visitor"
}
