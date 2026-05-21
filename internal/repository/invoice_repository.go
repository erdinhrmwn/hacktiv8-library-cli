package repository

import (
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

// Create membuat tagihan denda baru dengan status 'unpaid'.
func (r *InvoiceRepository) Create(userID, loanID int, amount float64) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO invoices (user_id, loan_id, amount, issue_date, status)
		VALUES (?, ?, ?, ?, 'unpaid')`,
		userID, loanID, amount, time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindByVisitorID mengambil semua invoice milik visitor tertentu.
func (r *InvoiceRepository) FindByVisitorID(visitorID int) ([]model.Invoice, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, loan_id, amount, issue_date, status
		FROM invoices
		WHERE user_id = ?`, visitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var inv model.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.UserID, &inv.LoanID,
			&inv.Amount, &inv.IssueDate, &inv.Status,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}
