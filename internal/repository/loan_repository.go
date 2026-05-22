package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type LoanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) GetAllLoans(ctx context.Context) ([]model.Loan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.visitor_id, l.staff_id, l.book_id,
		       l.borrow_date, l.due_date, l.return_date, l.status,
		       b.title, u.name
		FROM loans l
		JOIN books b ON b.id = l.book_id
		JOIN users u ON u.id = l.visitor_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		var book model.Book
		var visitor model.User
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status,
			&book.Title, &visitor.Name); err != nil {
			return nil, err
		}
		l.Book = &book
		l.Visitor = &visitor
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

func (r *LoanRepository) GetLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.visitor_id, l.staff_id, l.book_id,
		       l.borrow_date, l.due_date, l.return_date, l.status,
		       b.title, u.name
		FROM loans l
		JOIN books b ON b.id = l.book_id
		JOIN users u ON u.id = l.visitor_id
		WHERE l.id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var l model.Loan
	var book model.Book
	var visitor model.User
	if rows.Next() {
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status,
			&book.Title, &visitor.Name); err != nil {
			return nil, err
		}
		l.Book = &book
		l.Visitor = &visitor
	}
	return &l, rows.Err()
}

func (r *LoanRepository) CreateLoan(ctx context.Context, visitorID, staffID, bookID int) error {
	now := time.Now()
	dueDate := now.Add(7 * 24 * time.Hour)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO loans (visitor_id, staff_id, book_id, borrow_date, due_date, status) VALUES (?, ?, ?, ?, ?, 'active')`,
		visitorID, staffID, bookID, now, dueDate)
	return err
}

func (r *LoanRepository) ReturnLoan(ctx context.Context, loanID int) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE loans SET return_date = ?, status = 'returned' WHERE id = ?`,
		now, loanID)
	return err
}

func (r *LoanRepository) GetActiveLoansByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.visitor_id, l.staff_id, l.book_id,
		       l.borrow_date, l.due_date, l.return_date, l.status,
		       b.title
		FROM loans l
		JOIN books b ON b.id = l.book_id
		WHERE l.visitor_id = ? AND l.status = 'active'`, visitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		var book model.Book
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status,
			&book.Title); err != nil {
			return nil, err
		}
		l.Book = &book
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

type MostBorrowedBook struct {
	BookID      int
	Title       string
	BorrowCount int
}

func (r *LoanRepository) GetMostBorrowedBooks(ctx context.Context) ([]MostBorrowedBook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id, b.title, COUNT(l.id) as borrow_count
		FROM books b
		JOIN loans l ON l.book_id = b.id
		GROUP BY b.id, b.title
		ORDER BY borrow_count DESC
		LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MostBorrowedBook
	for rows.Next() {
		var m MostBorrowedBook
		if err := rows.Scan(&m.BookID, &m.Title, &m.BorrowCount); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, rows.Err()
}
