package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type InvoiceService struct {
	invoiceRepository *repository.InvoiceRepository
}

func NewInvoiceService(invoiceRepository *repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoiceRepository: invoiceRepository}
}

func (s *InvoiceService) GetUnpaidByUserID(ctx context.Context, userID int) ([]model.Invoice, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID user tidak valid")
	}
	return s.invoiceRepository.GetUnpaidInvoicesByUserID(ctx, userID)
}

func (s *InvoiceService) GetUnpaidInvoices(ctx context.Context) ([]model.Invoice, error) {
	return s.invoiceRepository.GetUnpaidInvoices(ctx)
}
