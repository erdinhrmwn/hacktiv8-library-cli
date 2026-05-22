package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) GetAll(ctx context.Context) ([]model.Invoice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM invoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var i model.Invoice
		if err := rows.Scan(&i.ID, &i.UserID, &i.LoanID, &i.Amount, &i.IssueDate, &i.Status); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, rows.Err()
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id int) (*model.Invoice, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM invoices WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var i model.Invoice
	if rows.Next() {
		if err := rows.Scan(&i.ID, &i.UserID, &i.LoanID, &i.Amount, &i.IssueDate, &i.Status); err != nil {
			return nil, err
		}
	}
	return &i, rows.Err()
}

func (r *InvoiceRepository) GetUnpaidByUserID(ctx context.Context, userID int) ([]model.Invoice, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT * FROM invoices WHERE user_id = ? AND status = 'unpaid'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var i model.Invoice
		if err := rows.Scan(&i.ID, &i.UserID, &i.LoanID, &i.Amount, &i.IssueDate, &i.Status); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, rows.Err()
}

func (r *InvoiceRepository) Create(ctx context.Context, userID, loanID int, amount float64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO invoices (user_id, loan_id, amount, issue_date, status) VALUES (?, ?, ?, ?, 'unpaid')`,
		userID, loanID, amount, time.Now())
	return err
}

func (r *InvoiceRepository) MarkPaid(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE invoices SET status = 'paid' WHERE id = ?`, id)
	return err
}
