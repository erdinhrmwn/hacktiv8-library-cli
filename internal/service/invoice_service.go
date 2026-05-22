package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type InvoiceService struct {
	invoiceRepository *repository.InvoiceRepository
}

func NewInvoiceService(invoiceRepository *repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoiceRepository: invoiceRepository}
}

func (s *InvoiceService) GetUnpaidByUserID(ctx context.Context, userID int) ([]repository.InvoiceDetail, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID user tidak valid")
	}
	return s.invoiceRepository.GetUnpaidInvoicesByUserID(ctx, userID)
}

func (s *InvoiceService) GetAllInvoices(ctx context.Context) ([]repository.InvoiceDetail, error) {
	return s.invoiceRepository.GetAllInvoices(ctx)
}
