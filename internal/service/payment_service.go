package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type PaymentService struct {
	paymentRepository *repository.PaymentRepository
	invoiceRepository *repository.InvoiceRepository
}

func NewPaymentService(paymentRepository *repository.PaymentRepository, invoiceRepository *repository.InvoiceRepository) *PaymentService {
	return &PaymentService{paymentRepository: paymentRepository, invoiceRepository: invoiceRepository}
}

func (s *PaymentService) PayInvoice(ctx context.Context, invoiceID int, amount float64, method string) error {
	if invoiceID <= 0 {
		return fmt.Errorf("ID invoice tidak valid")
	}
	if amount <= 0 {
		return fmt.Errorf("jumlah pembayaran tidak valid")
	}
	if method == "" {
		return fmt.Errorf("metode pembayaran tidak boleh kosong")
	}

	invoice, err := s.invoiceRepository.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}
	if invoice.ID == 0 {
		return fmt.Errorf("invoice tidak ditemukan")
	}
	if invoice.Status == "paid" {
		return fmt.Errorf("invoice sudah dibayar")
	}

	payment := &model.Payment{
		InvoiceID: invoiceID,
		Amount:    amount,
		Method:    method,
	}
	if err := s.paymentRepository.CreatePayment(ctx, payment); err != nil {
		return err
	}

	return s.invoiceRepository.MarkPaid(ctx, invoiceID)
}
