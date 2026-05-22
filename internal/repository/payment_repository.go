package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (invoice_id, amount, date, method) VALUES (?, ?, ?, ?)`,
		payment.InvoiceID, payment.Amount, time.Now(), payment.Method)
	return err
}
