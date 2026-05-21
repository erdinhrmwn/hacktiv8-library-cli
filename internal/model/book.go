package model

type Book struct {
	ID       int
	ISBN     string
	Title    string
	AuthorID int
	Genre    string
	Stock    int
	Author   *Author
}
