package controller

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type PaymentController struct {
	paymentService *service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (c *PaymentController) PayInvoice(ctx context.Context, invoiceID int, amount float64, method string) error {
	if invoiceID <= 0 {
		return fmt.Errorf("ID invoice tidak valid")
	}
	if amount <= 0 {
		return fmt.Errorf("jumlah pembayaran tidak valid")
	}
	if method == "" {
		return fmt.Errorf("metode pembayaran tidak boleh kosong")
	}

	return c.paymentService.PayInvoice(ctx, invoiceID, amount, method)
}
