package controller

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type InvoiceController struct {
	invoiceRepository *repository.InvoiceRepository
	loanRepository    *repository.LoanRepository
}

func NewInvoiceController(invoiceRepository *repository.InvoiceRepository, loanRepository *repository.LoanRepository) *InvoiceController {
	return &InvoiceController{invoiceRepository: invoiceRepository, loanRepository: loanRepository}
}

func (c *InvoiceController) GetUnpaidByUserID(ctx context.Context, userID int) ([]model.Invoice, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID user tidak valid")
	}
	return c.invoiceRepository.GetUnpaidInvoicesByUserID(ctx, userID)
}

func (c *InvoiceController) GetUnpaidInvoices(ctx context.Context) ([]model.Invoice, error) {
	return c.invoiceRepository.GetUnpaidInvoices(ctx)
}
