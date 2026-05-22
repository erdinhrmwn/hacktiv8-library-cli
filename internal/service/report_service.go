package service

import (
	"context"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type ReportService struct {
	loanRepository    *repository.LoanRepository
	invoiceRepository *repository.InvoiceRepository
}

func NewReportService(loanRepository *repository.LoanRepository, invoiceRepository *repository.InvoiceRepository) *ReportService {
	return &ReportService{loanRepository: loanRepository, invoiceRepository: invoiceRepository}
}

func (s *ReportService) GetMostBorrowedBooks(ctx context.Context) ([]repository.MostBorrowedBook, error) {
	return s.loanRepository.GetMostBorrowedBooks(ctx)
}

func (s *ReportService) GetTotalPaidFines(ctx context.Context) (float64, error) {
	return s.invoiceRepository.GetTotalPaidFines(ctx)
}

func (s *ReportService) GetTopUsersByFines(ctx context.Context) ([]repository.UserFine, error) {
	return s.invoiceRepository.GetTopUsersByFines(ctx)
}
