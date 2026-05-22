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

func (r *InvoiceRepository) GetAllInvoices(ctx context.Context) ([]model.Invoice, error) {
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

func (r *InvoiceRepository) GetInvoiceByID(ctx context.Context, id int) (*model.Invoice, error) {
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

func (r *InvoiceRepository) GetUnpaidInvoicesByUserID(ctx context.Context, userID int) ([]model.Invoice, error) {
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

func (r *InvoiceRepository) CreateInvoice(ctx context.Context, userID, loanID int, amount float64) error {
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

func (r *InvoiceRepository) GetUnpaidInvoices(ctx context.Context) ([]model.Invoice, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM invoices WHERE status = 'unpaid'`)
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

func (r *InvoiceRepository) GetTotalPaidFines(ctx context.Context) (float64, error) {
	var total sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT SUM(amount) FROM invoices WHERE status = 'paid'`).Scan(&total)
	if err != nil {
		return 0, err
	}
	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

type UserFine struct {
	UserID    int
	Name      string
	TotalFine float64
}

func (r *InvoiceRepository) GetTopUsersByFines(ctx context.Context) ([]UserFine, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.name, COALESCE(SUM(i.amount), 0) as total_fine
		FROM users u
		JOIN invoices i ON i.user_id = u.id
		WHERE i.status = 'paid'
		GROUP BY u.id, u.name
		ORDER BY total_fine DESC
		LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserFine
	for rows.Next() {
		var uf UserFine
		if err := rows.Scan(&uf.UserID, &uf.Name, &uf.TotalFine); err != nil {
			return nil, err
		}
		results = append(results, uf)
	}
	return results, rows.Err()
}
