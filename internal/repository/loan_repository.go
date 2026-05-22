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
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM loans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

func (r *LoanRepository) GetLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM loans WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var l model.Loan
	if rows.Next() {
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status); err != nil {
			return nil, err
		}
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
	rows, err := r.db.QueryContext(ctx,
		`SELECT * FROM loans WHERE visitor_id = ? AND status = 'active'`, visitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		if err := rows.Scan(&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.ReturnDate, &l.Status); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}
