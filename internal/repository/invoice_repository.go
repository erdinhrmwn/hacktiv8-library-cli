package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type InvoiceDetail struct {
	ID        int
	UserName  string
	BookTitle string
	Amount    float64
	Status    string
	IssueDate time.Time
}

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) GetAllInvoices(ctx context.Context) ([]InvoiceDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, u.name, b.title, i.amount, i.status, i.issue_date
		FROM invoices i
		JOIN users u ON u.id = i.user_id
		JOIN loans l ON l.id = i.loan_id
		JOIN books b ON b.id = l.book_id
		ORDER BY i.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []InvoiceDetail
	for rows.Next() {
		var d InvoiceDetail
		if err := rows.Scan(&d.ID, &d.UserName, &d.BookTitle, &d.Amount, &d.Status, &d.IssueDate); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
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

func (r *InvoiceRepository) GetInvoicesByUserID(ctx context.Context, userID int) ([]InvoiceDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, u.name, b.title, i.amount, i.status, i.issue_date
		FROM invoices i
		JOIN users u ON u.id = i.user_id
		JOIN loans l ON l.id = i.loan_id
		JOIN books b ON b.id = l.book_id
		WHERE i.user_id = ?
		ORDER BY i.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []InvoiceDetail
	for rows.Next() {
		var d InvoiceDetail
		if err := rows.Scan(&d.ID, &d.UserName, &d.BookTitle, &d.Amount, &d.Status, &d.IssueDate); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *InvoiceRepository) GetUnpaidInvoicesByUserID(ctx context.Context, userID int) ([]InvoiceDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, u.name, b.title, i.amount, i.status, i.issue_date
		FROM invoices i
		JOIN users u ON u.id = i.user_id
		JOIN loans l ON l.id = i.loan_id
		JOIN books b ON b.id = l.book_id
		WHERE i.user_id = ? AND i.status = 'unpaid'
		ORDER BY i.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []InvoiceDetail
	for rows.Next() {
		var d InvoiceDetail
		if err := rows.Scan(&d.ID, &d.UserName, &d.BookTitle, &d.Amount, &d.Status, &d.IssueDate); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
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
