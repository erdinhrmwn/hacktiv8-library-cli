package controller

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type InvoiceController struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceController(invoiceService *service.InvoiceService) *InvoiceController {
	return &InvoiceController{invoiceService: invoiceService}
}

func (c *InvoiceController) GetUnpaidByUserID(ctx context.Context, userID int) ([]model.Invoice, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID user tidak valid")
	}
	return c.invoiceService.GetUnpaidByUserID(ctx, userID)
}

func (c *InvoiceController) GetUnpaidInvoices(ctx context.Context) ([]model.Invoice, error) {
	return c.invoiceService.GetUnpaidInvoices(ctx)
}

func (c *InvoiceController) GetAllInvoices(ctx context.Context) ([]model.Invoice, error) {
	return c.invoiceService.GetAllInvoices(ctx)
}
